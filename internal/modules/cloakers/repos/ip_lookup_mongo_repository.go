package repos

import (
	"context"
	"errors"
	"time"

	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/models"
	"github.com/italoservio/serviosoftware_ads/pkg/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoIPLookupRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoIPLookupRepository(d *db.DB) *MongoIPLookupRepository {
	database := d.Client.Database("ads")
	collection := database.Collection("ip_lookups")

	repo := MongoIPLookupRepository{db: database, coll: collection}
	repo.createIndices()

	return &repo
}

func (r *MongoIPLookupRepository) createIndices() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Índice único parcial: ipPattern é único apenas quando deletedAt é null
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "ipPattern", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{"deletedAt": bson.M{"$eq": nil}}),
	})

	if err != nil {
		panic(err)
	}
}

func (r *MongoIPLookupRepository) GetByIPPattern(pattern string) (*models.IPLookup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ipLookup models.IPLookup

	filter := bson.M{"ipPattern": pattern, "deletedAt": bson.M{"$eq": nil}}
	err := r.coll.FindOne(ctx, filter).Decode(&ipLookup)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &ipLookup, nil
}

func (r *MongoIPLookupRepository) Create(ipLookup *models.IPLookup) (*models.IPLookup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ipLookup.CreatedAt = time.Now()
	ipLookup.UpdatedAt = time.Now()
	ipLookup.AccessCount = 0

	inserted, err := r.coll.InsertOne(ctx, ipLookup)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return r.GetByIPPattern(ipLookup.IPPattern)
		}

		return nil, err
	}

	ipLookup.ID = inserted.InsertedID.(bson.ObjectID)
	return ipLookup, nil
}

func (r *MongoIPLookupRepository) IncrementAccessCount(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID, "deletedAt": bson.M{"$eq": nil}}
	update := bson.M{
		"$inc": bson.M{"accessCount": 1},
		"$set": bson.M{"updatedAt": time.Now()},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("ip lookup nao encontrado")
	}

	return nil
}
