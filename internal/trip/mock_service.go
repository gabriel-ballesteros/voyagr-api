package trip

import (
	"context"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
	"github.com/google/uuid"
)

type MockService interface {
	GetAll(ctx context.Context, user_id string) ([]domain.Trip, error)
	Get(ctx context.Context, id string) (domain.Trip, error)
	Store(ctx context.Context, name string, description string, start string, end string, owner string, sharedWith []string, itinerary []domain.ItineraryElement) (domain.Trip, error)
	Update(ctx context.Context, id string, name string, description string, start string, end string, owner string, sharedWith []string, itinerary []domain.ItineraryElement) (domain.Trip, error)
	Delete(ctx context.Context, id string) error
}

type mockService struct {
	db *map[string]domain.Trip
}

func NewMockService(db *map[string]domain.Trip) MockService {
	return &mockService{db: db}
}

func (s *mockService) GetAll(ctx context.Context, user_id string) ([]domain.Trip, error) {
	var tripList []domain.Trip
	for _, trip := range *s.db {
		if trip.Owner == user_id {
			tripList = append(tripList, trip)
		}
	}
	if len(tripList) == 0 {
		return nil, web.NewError(404, "There are no trips for this user")
	}
	return tripList, nil
}

func (s *mockService) Get(ctx context.Context, id string) (domain.Trip, error) {
	if _, exists := (*s.db)[id]; !exists {
		return domain.Trip{}, web.NewError(404, "The trip with id "+id+" does not exist")
	}
	return (*s.db)[id], nil
}

func (s *mockService) Store(ctx context.Context, name string, description string, start string, end string, owner string, sharedWith []string, itinerary []domain.ItineraryElement) (domain.Trip, error) {
	id := uuid.New()
	newTrip := domain.Trip{
		ID:          id.String(),
		Name:        name,
		Description: description,
		Start:       start,
		End:         end,
		Owner:       owner,
		SharedWith:  sharedWith,
		Itinerary:   itinerary,
	}
	(*s.db)[id.String()] = newTrip
	return newTrip, nil
}
func (s *mockService) Update(ctx context.Context, id string, name string, description string, start string, end string, owner string, sharedWith []string, itinerary []domain.ItineraryElement) (domain.Trip, error) {
	_, err := s.Get(ctx, id)
	if err != nil {
		return domain.Trip{}, web.NewError(404, err.Error())
	}

	updatedTrip := domain.Trip{
		Name:        name,
		Description: description,
		Start:       start,
		End:         end,
		Owner:       owner,
		SharedWith:  sharedWith,
		Itinerary:   itinerary,
	}

	(*s.db)[id] = updatedTrip
	return updatedTrip, nil
}
func (s *mockService) Delete(ctx context.Context, id string) error {

	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	delete(*s.db, id)
	return nil
}
