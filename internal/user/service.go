package user

import (
	"context"
	"fmt"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/utils"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	Get(ctx context.Context, email string) (domain.User, error)
	Store(ctx context.Context, name string, email string) (domain.User, error)
	Update(ctx context.Context, email string, name string) (domain.User, error)
	ResetPassword(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) error
	Delete(ctx context.Context, email string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{
		repository: r,
	}
}

// Get function: get a single user by email, returns 404 if not found
func (s *service) Get(ctx context.Context, email string) (domain.User, error) {
	u, err := s.repository.Get(ctx, email)
	if err != nil {
		errMessage := fmt.Sprintf("The user with email %s does not exist", email)
		return domain.User{}, web.NewError(404, errMessage)
	} else {
		return u, nil
	}
}

// Store function, creates a user, returns 409 if user is already in db or 500 if has any database error
func (s *service) Store(ctx context.Context, name string, email string) (domain.User, error) {
	_, err := s.repository.Get(ctx, email)
	if err != mongo.ErrNoDocuments {
		return domain.User{}, web.NewErrorf(409, "User already in database")
	}
	password := utils.GenerateRandomString(12)
	var newUser domain.User = domain.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	resultUser, storeErr := s.repository.Save(ctx, newUser)

	if storeErr != nil {
		return domain.User{}, web.NewErrorf(500, storeErr.Error())
	}

	return resultUser, nil
}

// Update function, searches a user by email and updates the fields
// If the user is not found, it returns 404
// else, it updates the fields and returns 500 in case of error while updating
func (s *service) Update(ctx context.Context, email string, name string) (domain.User, error) {

	userToUpdate, err := s.Get(ctx, email)
	if err != nil {
		return domain.User{}, web.NewError(404, err.Error())
	}
	userToUpdate.Name = name

	if err := s.repository.Update(ctx, userToUpdate); err != nil {
		return domain.User{}, web.NewError(500, err.Error())
	}

	return userToUpdate, nil
}

// the ResetPassword function hard resets the password to a random 12 alphanumeric string
func (s *service) ResetPassword(ctx context.Context, email string) error {
	password := utils.GenerateRandomString(12)
	if err := s.repository.SetPassword(ctx, email, password); err != nil {
		return err
	}
	return nil
}

// Change password function: searches for a user by email
// If the user exists, it checks its current password and compares it against the input
// If the old passwords match, it tries to update the password, else it returns 401
// If the password update returns error, it returns 500
func (s *service) ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) error {
	u, err := s.repository.Get(ctx, email)
	if err != nil {
		errMessage := fmt.Sprintf("El user con email %s no existe en la base de datos", email)
		return web.NewError(404, errMessage)
	} else {
		if u.Password == oldPassword {
			if err := s.repository.SetPassword(ctx, email, newPassword); err != nil {
				return web.NewError(500, err.Error())
			}
		} else {
			return web.NewError(401, "Wrong user and/or password")
		}
	}
	return nil
}

// Delete function: searches for a user by email and deletes it
// Returns 404 if the user is not found
func (s *service) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)

	if err != nil {
		return web.NewError(404, err.Error())
	}

	return nil
}
