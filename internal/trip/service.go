package trip

import (
	"context"
	"fmt"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
)

type Service interface {
	GetAll(ctx context.Context, user_id int) ([]domain.Trip, error)
	Get(ctx context.Context, id int) (domain.Trip, error)
	Store(ctx context.Context, name string, start string, end string, owner int, sharedWith []int, itinerary []interface{}) (domain.Trip, error)
	Update(ctx context.Context, id int, name string, start string, end string, owner int, sharedWith []int, itinerary []interface{}) (domain.Trip, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{
		repository: r,
	}
}

// Funcion getall: obtiene todos los trips
// si el slice devuelto por el repository esta vacio, devuelve 404
// si devuelve otro error, devuelve la descripcion con codigo internal error
// else, devuelve un slice de trips y error nulo
func (s *service) GetAll(ctx context.Context, user_id int) ([]domain.Trip, error) {
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

// Funcion get: obtiene un trip buscando por id
// si el get de repositorio devuelve un error, se asume que ese trip no existe devolviendo 404
// else, retorna el trip encontrado y error nulo
func (s *service) Get(ctx context.Context, id int) (domain.Trip, error) {
	wh, err := s.repository.Get(ctx, id)
	if err != nil {
		errMessage := fmt.Sprintf("El trip con id %d no existe en la base de datos", id)
		return domain.Trip{}, web.NewError(404, errMessage)
	} else {
		return wh, nil
	}
}

// Funcion store: da de alta un trip en la bd
// si la base da un error, devuelve status 500
// else, lo da de alta
func (s *service) Store(ctx context.Context, name string, start string, end string, owner int, sharedWith []int, itinerary []interface{}) (domain.Trip, error) {

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

// Funcion update: busca un trip por id, actualiza los campos del objeto, luego lo actualiza en la bd
// Si no se encuentra el id del trip, devuelve 404
// else, actualiza
func (s *service) Update(ctx context.Context, id int, name string, start string, end string, owner int, sharedWith []int, itinerary []interface{}) (domain.Trip, error) {

	tripToUpdate, err := s.Get(ctx, id)
	if err != nil {
		return domain.Trip{}, web.NewError(404, err.Error())
	}

	tripToUpdate.ID = id
	tripToUpdate.Name = name
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

// funcion delete: busca trip por id y lo borra
// si no se encuentra el id en la base, devuelve 404, else borra
func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)

	if err != nil {
		return web.NewError(404, err.Error())
	}

	return nil
}
