package mongo

import (
	"context"
	"errors"
	"film-list/internal/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log/slog"
	"time"
)

type Storage struct {
	client *mongo.Client
}

func New() (*Storage, error) {
	mongoConnectionCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(mongoConnectionCtx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		slog.Error("failed to connect to MongoDB", err)
	}

	err = client.Ping(mongoConnectionCtx, nil)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("connected to MongoDB!")

	return &Storage{client: client}, nil
}

func (s *Storage) GetFilms() ([]dto.Film, error) {
	var filmsSlice []dto.Film

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("film-list").Collection("films")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		slog.Error("failed to find films", err)

		return nil, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			slog.Error("failed to close cursor", err)
		}
	}()

	for cursor.Next(ctx) {
		var film dto.Film
		err := cursor.Decode(&film)
		if err != nil {
			slog.Error("failed to decode film", err)
			// TODO: Вернуть + обработать ошибку
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
		slog.Error("failed to insert film", err)

		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		slog.Error("failed to convert inserted ID to ObjectID")

		return "", errors.New("failed to convert inserted ID to ObjectID")
	}

	return oid.Hex(), nil
}

func (s *Storage) DeleteFilm(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("film-list").Collection("films")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("failed to convert id to ObjectID:", "error", err)

		return 0, err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		slog.Error("failed to delete film:", "id", objectId, "error", err)

		return 0, err
	}

	if result.DeletedCount == 0 {
		slog.Warn("no films were deleted:", "id", objectId)

		return 0, nil
	}

	slog.Info("deleted film successfully:", "id", objectId)

	return result.DeletedCount, nil
}
