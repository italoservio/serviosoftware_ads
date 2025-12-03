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

type MongoCloakerRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoCloakerRepository(d *db.DB) *MongoCloakerRepository {
	database := d.Client.Database("ads")
	collection := database.Collection("cloakers")

	repo := MongoCloakerRepository{db: database, coll: collection}
	repo.createIndices()

	return &repo
}

func (r *MongoCloakerRepository) createIndices() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "userId", Value: 1}},
	})

	if err != nil {
		panic(err)
	}

	_, err = r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		panic(err)
	}
}

func (r *MongoCloakerRepository) GenerateID() string {
	return bson.NewObjectID().Hex()
}

func (r *MongoCloakerRepository) ToID(id string) bson.ObjectID {
	objectID, _ := bson.ObjectIDFromHex(id)
	return objectID
}

func (r *MongoCloakerRepository) CreateWithID(id string, cloaker *models.Cloaker) (*models.Cloaker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cloaker.ID = objectID
	cloaker.CreatedAt = time.Now()
	cloaker.UpdatedAt = time.Now()
	cloaker.DeletedAt = nil
	cloaker.IsActive = false

	_, err = r.coll.InsertOne(ctx, cloaker)
	if err != nil {
		return nil, err
	}

	return cloaker, nil
}

func (r *MongoCloakerRepository) Create(cloaker *models.Cloaker) (*models.Cloaker, error) {
	id := r.GenerateID()
	return r.CreateWithID(id, cloaker)
}

func (r *MongoCloakerRepository) GetByID(id string) (*models.Cloaker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cloaker models.Cloaker

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var filter = bson.M{"_id": objectID, "deletedAt": nil}
	err = r.coll.FindOne(ctx, filter).Decode(&cloaker)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &cloaker, nil
}

func (r *MongoCloakerRepository) List(input *ListCloakersInput) (*ListCloakersOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	match := bson.M{"deletedAt": bson.M{"$eq": nil}}

	if input.UserID != nil {
		userID, err := bson.ObjectIDFromHex(*input.UserID)
		if err != nil {
			return nil, err
		}
		match["userId"] = userID
	}

	sortBy := *input.SortBy

	sortOrder := 1
	if *input.Order == "desc" {
		sortOrder = -1
	}

	cursor, err := r.coll.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$facet", Value: bson.M{
			"total": []bson.M{{"$count": "count"}},
			"items": []bson.M{
				{"$sort": bson.M{sortBy: sortOrder}},
				{"$skip": (input.Page - 1) * input.Limit},
				{"$limit": input.Limit},
			}},
		}},
		{{Key: "$project", Value: bson.M{
			"total": bson.M{
				"$arrayElemAt": []interface{}{"$total.count", 0},
			},
			"items": "$items",
		}}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []ListCloakersOutput
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &ListCloakersOutput{}, nil
	}

	output := results[0]
	return &output, nil
}

func (r *MongoCloakerRepository) UpdateByID(id string, cloaker *models.Cloaker) (*models.Cloaker, error) {
	cloakerID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cloaker.ID = cloakerID

	updateDoc := bson.M{}

	if cloaker.URL != "" {
		updateDoc["url"] = cloaker.URL
	}

	if cloaker.WhiteURL != "" {
		updateDoc["whiteUrl"] = cloaker.WhiteURL
	}

	if cloaker.BlackURL != "" {
		updateDoc["blackUrl"] = cloaker.BlackURL
	}

	if cloaker.Config.AllowOnlyMobile || !cloaker.Config.AllowOnlyMobile {
		updateDoc["config"] = cloaker.Config
	}

	updateDoc["isActive"] = cloaker.IsActive

	if len(updateDoc) == 0 {
		return nil, errors.New("nao ha campos para atualizar")
	}

	updateDoc["updatedAt"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": cloaker.ID, "deletedAt": nil}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": updateDoc}, opts).Decode(&cloaker)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("cloaker nao encontrado")
		}

		return nil, err
	}

	return cloaker, err
}

func (r *MongoCloakerRepository) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.coll.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
