package trip

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
)

// Repository encapsulates the storage of a trip.
type Repository interface {
	GetAll(ctx context.Context, user_id string) ([]domain.Trip, error)
	Get(ctx context.Context, id string) (domain.Trip, error)
	Save(ctx context.Context, t domain.Trip) (domain.Trip, error)
	Update(ctx context.Context, w domain.Trip) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *mongo.Collection
}

func NewRepository(db *mongo.Collection) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, user_id string) ([]domain.Trip, error) {
	filter := bson.M{"owner": user_id}
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return []domain.Trip{}, err
	}
	var results []domain.Trip
	fmt.Println(user_id)
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Println(err)
		return []domain.Trip{}, err
	}
	return results, nil
}

func (r *repository) Get(ctx context.Context, id string) (domain.Trip, error) {
	var resultTrip domain.Trip
	objID, _ := primitive.ObjectIDFromHex(id)
	err := r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&resultTrip)
	if err != nil {
		fmt.Println(err)
		return domain.Trip{}, err
	}

	return resultTrip, nil
}

func (r *repository) Save(ctx context.Context, t domain.Trip) (domain.Trip, error) {
	var resultTrip domain.Trip
	insertResult, err := r.db.InsertOne(context.TODO(), t)
	if err != nil {
		return domain.Trip{}, err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	err = r.db.FindOne(ctx, bson.M{"_id": insertResult.InsertedID.(primitive.ObjectID)}).Decode(&resultTrip)

	return resultTrip, nil
}

func (r *repository) Update(ctx context.Context, updatedTrip domain.Trip) error {
	objID, _ := primitive.ObjectIDFromHex(updatedTrip.ID)

	// Not the best way to do this, but it works and we're only editing a transient object.
	updatedTrip.ID = ""
	update := bson.D{
		{"$set", updatedTrip},
	}

	filter := bson.D{{"_id", objID}}
	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	deleteResult, err := r.db.DeleteOne(ctx, bson.D{{"_id", objID}})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil
}
