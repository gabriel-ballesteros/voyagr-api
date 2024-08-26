package trip

import (
	"context"
	"fmt"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
)

type Service interface {
	GetAll(ctx context.Context, user_id string) ([]domain.Trip, error)
	Get(ctx context.Context, id string) (domain.Trip, error)
	Store(ctx context.Context, name string, description string, start string, end string, owner string, sharedWith []int, itinerary []domain.ItineraryElement) (domain.Trip, error)
	Update(ctx context.Context, id string, name string, description string, start string, end string, owner string, sharedWith []int, itinerary []domain.ItineraryElement) (domain.Trip, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{
		repository: r,
	}
}

// GetAll function: gets all trips from a single user_id, returns 500 if has any error
func (s *service) GetAll(ctx context.Context, user_id string) ([]domain.Trip, error) {
	trips, err := s.repository.GetAll(ctx, user_id)
	if err == nil && len(trips) == 0 {
		return nil, web.NewError(404, "No existen trips dados de alta en la base de datos")
	} else if err != nil {
		fmt.Println(err)
		return nil, web.NewError(500, "Error en la base de datos")
	} else {
		return trips, nil
	}
}

// Get function: get a single trip by id, returns 404 if not found
func (s *service) Get(ctx context.Context, id string) (domain.Trip, error) {
	wh, err := s.repository.Get(ctx, id)
	if err != nil {
		errMessage := fmt.Sprintf("El trip con id %s no existe en la base de datos", id)
		return domain.Trip{}, web.NewError(404, errMessage)
	} else {
		return wh, nil
	}
}

// Store function, creates a trip, returns 500 if has any error
func (s *service) Store(ctx context.Context, name string, description string,
	start string, end string, owner string, sharedWith []int, itinerary []domain.ItineraryElement) (domain.Trip, error) {

	var newTrip domain.Trip = domain.Trip{
		//ID:         id,
		Name:       name,
		Start:      start,
		End:        end,
		Owner:      owner,
		SharedWith: sharedWith,
		Itinerary:  itinerary,
	}

	storeErr := s.repository.Save(ctx, newTrip)

	if storeErr != nil {
		return domain.Trip{}, web.NewErrorf(409, storeErr.Error())
	}

	return newTrip, nil
}

// Update function, searches a trip by id and updates the fields
// If the trip is not found, it returns 404
// else, it updates the fields
func (s *service) Update(ctx context.Context, id string, name string, description string,
	start string, end string, owner string, sharedWith []int, itinerary []domain.ItineraryElement) (domain.Trip, error) {

	tripToUpdate, err := s.Get(ctx, id)
	if err != nil {
		return domain.Trip{}, web.NewError(404, err.Error())
	}

	tripToUpdate.ID = id
	tripToUpdate.Name = name
	tripToUpdate.Description = description
	tripToUpdate.Start = start
	tripToUpdate.End = end
	tripToUpdate.Owner = owner
	tripToUpdate.SharedWith = sharedWith
	tripToUpdate.Itinerary = itinerary

	if err := s.repository.Update(ctx, tripToUpdate); err != nil {
		return domain.Trip{}, web.NewError(404, err.Error())
	}

	return tripToUpdate, nil
}

// Delete function: searchesa  trip by id and deletes it
// Returns 404 if the trip is not found
func (s *service) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)

	if err != nil {
		return web.NewError(404, err.Error())
	}

	return nil
}
