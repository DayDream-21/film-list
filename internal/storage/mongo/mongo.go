package mongo

import (
	"context"
	"errors"
	"film-list/internal/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

type Storage struct {
	client *mongo.Client
	log    *slog.Logger
}

func New(log *slog.Logger) (*Storage, error) {
	mongoConnectionCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(mongoConnectionCtx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		log.Error("failed to connect to MongoDB", err)

		return nil, err
	}

	err = client.Ping(mongoConnectionCtx, nil)
	if err != nil {
		log.Error("failed to ping MongoDB", err)
	}

	log.Info("connected to MongoDB!")

	return &Storage{client: client, log: log}, nil
}

func (s *Storage) GetFilms() ([]dto.Film, error) {
	var filmsSlice []dto.Film

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: вынести имя БД и коллекцию в конфиг
	collection := s.client.Database("film-list").Collection("films")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error("failed to find films:", "error", err)

		return nil, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			s.log.Error("failed to close cursor:", "error", err)
		}
	}()

	for cursor.Next(ctx) {
		var film dto.Film
		err := cursor.Decode(&film)
		if err != nil {
			s.log.Error("failed to decode film:", "error", err)

			return nil, err
		}
		filmsSlice = append(filmsSlice, film)
	}

	return filmsSlice, nil
}

func (s *Storage) SaveFilm(film dto.Film) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("film-list").Collection("films")

	result, err := collection.InsertOne(ctx, film)
	if err != nil {
		s.log.Error("failed to insert film:", "error", err)

		return "", err
	}

	s.log.Info("saved film successfully:", "id", result.InsertedID)

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err := errors.New("failed to convert inserted ID to ObjectID")

		s.log.Error("failed to convert inserted ID to ObjectID:", "error", err)

		return "", err
	}

	return oid.Hex(), nil
}

func (s *Storage) DeleteFilm(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("film-list").Collection("films")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.log.Error("failed to convert id to ObjectID:", "error", err)

		return 0, err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		s.log.Error("failed to delete film:", "id", objectId, "error", err)

		return 0, err
	}

	if result.DeletedCount == 0 {
		s.log.Warn("no films were deleted:", "id", objectId)

		return 0, nil
	}

	s.log.Info("deleted film successfully:", "id", objectId)

	return result.DeletedCount, nil
}
