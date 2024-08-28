package user

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
	Get(ctx context.Context, email string) (domain.User, error)
	Save(ctx context.Context, t domain.User) (domain.User, error)
	Update(ctx context.Context, w domain.User) error
	ResetPassword(ctx context.Context, email string, newPassword string) error
}

type repository struct {
	db *mongo.Collection
}

func NewRepository(db *mongo.Collection) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Get(ctx context.Context, email string) (domain.User, error) {
	var resultUser domain.User
	err := r.db.FindOne(ctx, bson.M{"email": email}).Decode(&resultUser)
	if err != nil {
		fmt.Println(err)
		return domain.User{}, err
	}

	return resultUser, nil
}

func (r *repository) Save(ctx context.Context, t domain.User) (domain.User, error) {
	var resultUser domain.User
	insertResult, err := r.db.InsertOne(ctx, t)
	if err != nil {
		return domain.User{}, err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	err = r.db.FindOne(ctx, bson.M{"_id": insertResult.InsertedID.(primitive.ObjectID)}).Decode(&resultUser)
	if err != nil {
		return domain.User{}, err
	}
	return resultUser, nil
}

func (r *repository) Update(ctx context.Context, updatedUser domain.User) error {

	// Not the best way to do this, but it works and we're only editing a transient object.
	update := bson.D{{"$set", updatedUser}}
	filter := bson.D{{"email", updatedUser.Email}}
	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) ResetPassword(ctx context.Context, email string, newPassword string) error {
	var resultUser domain.User
	filter := bson.M{"email": email}
	update := bson.D{{"$set", bson.D{{"password", newPassword}}}}
	err := r.db.FindOne(ctx, filter).Decode(&resultUser)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		_, err := r.db.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}
