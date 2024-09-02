package user

import (
	"context"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/utils"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
)

type MockService interface {
	Get(ctx context.Context, email string) (domain.User, error)
	Store(ctx context.Context, name string, email string) (domain.User, error)
	Update(ctx context.Context, email string, name string) (domain.User, error)
	ResetPassword(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) error
	Delete(ctx context.Context, email string) error
}

type mockService struct {
	db *map[string]domain.User
}

func NewMockService(db *map[string]domain.User) MockService {
	return &mockService{db: db}
}

func (s *mockService) Get(ctx context.Context, email string) (domain.User, error) {
	if _, exists := (*s.db)[email]; !exists {
		return domain.User{}, web.NewError(404, "The user with email "+email+" does not exist")
	}
	return (*s.db)[email], nil
}

func (s *mockService) Store(ctx context.Context, email string, name string) (domain.User, error) {
	_, err := s.Get(ctx, email)
	if err == nil {
		return domain.User{}, web.NewError(409, "An user with the email "+email+" already exists")
	}
	password := utils.GenerateRandomString(12)
	newUser := domain.User{
		Email:    email,
		Name:     name,
		Password: password,
	}
	(*s.db)[email] = newUser
	return newUser, nil
}
func (s *mockService) Update(ctx context.Context, email string, name string) (domain.User, error) {
	oldUser, err := s.Get(ctx, email)
	if err != nil {
		return domain.User{}, web.NewError(404, err.Error())
	}

	updatedUser := domain.User{
		Email:    email,
		Name:     name,
		Password: oldUser.Password,
	}

	(*s.db)[email] = updatedUser
	return updatedUser, nil
}
func (s *mockService) ResetPassword(ctx context.Context, email string) error {
	_, err := s.Get(ctx, email)
	if err != nil {
		return web.NewError(404, err.Error())
	}
	if email == "" {
		return web.NewError(500, "Internal server error")
	}
	return nil
}
func (s *mockService) ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) error {
	user, err := s.Get(ctx, email)
	if err != nil {
		return web.NewError(404, err.Error())
	}
	if email == "" {
		return web.NewError(500, "Internal server error")
	}
	if user.Password != oldPassword {
		return web.NewError(401, "Wrong user and/or password")
	}
	return nil
}
func (s *mockService) Delete(ctx context.Context, email string) error {

	_, err := s.Get(ctx, email)
	if err != nil {
		return err
	}

	delete(*s.db, email)
	return nil
}
