package myMongo

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
		slog.Error("Failed to connect to MongoDB", err)
	}

	err = client.Ping(mongoConnectionCtx, nil)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Connected to MongoDB!")

	return &Storage{client: client}, nil
}

func (s *Storage) GetFilms() ([]dto.Film, error) {
	var filmsSlice []dto.Film

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("film-list").Collection("films")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		slog.Error("Failed to find films", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var film dto.Film
		err := cursor.Decode(&film)
		if err != nil {
			slog.Error("Failed to decode film", err)
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
		slog.Error("Failed to insert film", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("Failed to convert inserted ID to ObjectID")
	}

	return oid.Hex(), nil
}

func (s *Storage) DeleteFilm(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("film-list").Collection("films")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("Failed to convert id to ObjectID:", "error", err)
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		slog.Error("Failed to delete film:", "id", objectId, "error", err)
	}

	if result.DeletedCount == 0 {
		slog.Warn("No films were deleted:", "id", objectId)
	}

	slog.Info("Deleted film successfully:", "id", objectId)

	return result.DeletedCount, nil
}
