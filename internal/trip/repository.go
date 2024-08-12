package trip

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
)

// Repository encapsulates the storage of a trip.
type Repository interface {
	GetAll(ctx context.Context, user_id int) ([]domain.Trip, error)
	Get(ctx context.Context, id int) (domain.Trip, error)
	Save(ctx context.Context, w domain.Trip) error
	Update(ctx context.Context, w domain.Trip) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *mongo.Collection
}

func NewRepository(db *mongo.Collection) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context, user_id int) ([]domain.Trip, error) {
	filter := bson.D{{"owner", user_id}}
	cursor, err := r.db.Find(context.TODO(), filter)
	if err != nil {
		return []domain.Trip{}, err
	}
	var results []domain.Trip
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []domain.Trip{}, err
	}
	return results, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Trip, error) {
	filter := bson.D{{"id", id}}
	var result domain.Trip

	err := r.db.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return domain.Trip{}, err
	}
	return result, nil
}

func (r *repository) Save(ctx context.Context, t domain.Trip) error {

	insertResult, err := r.db.InsertOne(context.TODO(), t)
	if err != nil {
		return err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}

func (r *repository) Update(ctx context.Context, updatedTrip domain.Trip) error {

	filter := bson.D{{"id", updatedTrip.ID}}
	_, err := r.db.UpdateOne(context.TODO(), filter, updatedTrip)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	deleteResult, err := r.db.DeleteOne(context.TODO(), bson.D{{"id", id}})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil
}
