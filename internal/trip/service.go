package trip

import (
	"cmp"
	"context"
	"fmt"
	"sort"

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
		return nil, web.NewError(404, "There are no trips for this user")
	} else if err != nil {
		fmt.Println(err)
		return nil, web.NewError(500, "Unexpected error with MongoDB")
	} else {
		return trips, nil
	}
}

// Get function: get a single trip by id, returns 404 if not found
func (s *service) Get(ctx context.Context, id string) (domain.Trip, error) {
	t, err := s.repository.Get(ctx, id)
	if err != nil {
		errMessage := fmt.Sprintf("The trip with id %s does not exist", id)
		return domain.Trip{}, web.NewError(404, errMessage)
	} else {
		return t, nil
	}
}

// Store function, creates a trip, returns 500 if has any error
func (s *service) Store(ctx context.Context, name string, description string,
	start string, end string, owner string, sharedWith []int, itinerary []domain.ItineraryElement) (domain.Trip, error) {

	var newTrip domain.Trip = domain.Trip{
		//ID:         id,
		Name:        name,
		Description: description,
		Start:       start,
		End:         end,
		Owner:       owner,
		SharedWith:  sharedWith,
		Itinerary:   itinerary,
	}

	resultTrip, storeErr := s.repository.Save(ctx, newTrip)

	if storeErr != nil {
		return domain.Trip{}, web.NewErrorf(409, storeErr.Error())
	}

	return resultTrip, nil
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

	// sorting the itinerary list by datetime before saving the updated data
	sort.Slice(itinerary, func(i, j int) bool {
		return cmp.Or(itinerary[i].Departure, itinerary[i].CheckIn, itinerary[i].EventDatetime) < cmp.Or(itinerary[j].Departure, itinerary[j].CheckIn, itinerary[i].EventDatetime)
	})

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
