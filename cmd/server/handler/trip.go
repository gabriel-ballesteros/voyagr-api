package handler

import (
	"strconv"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	trip "github.com/gabriel-ballesteros/voyagr-api/internal/trip"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
	"github.com/gin-gonic/gin"
)

type Trip struct {
	tripService trip.Service
}

func NewTrip(t trip.Service) *Trip {
	return &Trip{
		tripService: t,
	}
}

func (t *Trip) GetAll() gin.HandlerFunc {
	type response struct {
		Data []domain.Trip `json:"data"`
	}

	return func(c *gin.Context) {
		user_id := c.Query("user_id")
		wrhs, err := t.tripService.GetAll(c, user_id)

		if err != nil && wrhs == nil {
			code, _ := strconv.Atoi(err.Error()[0:3])
			c.JSON(code, gin.H{
				"error": err,
			})
			return
		} else {
			res := response{
				Data: wrhs,
			}
			c.JSON(200, res)
			return
		}
	}
}

func (t *Trip) Get() gin.HandlerFunc {
	type response struct {
		Data domain.Trip `json:"data"`
	}

	return func(c *gin.Context) {
		id := c.Param("id")

		tr, err := t.tripService.Get(c, id)
		if err != nil {
			c.JSON(404, web.NewError(404, "Trip with id "+id+" not found"))
			return
		}

		res := response{
			Data: tr,
		}

		c.JSON(200, res)
		return
	}
}

func (t *Trip) Store() gin.HandlerFunc {
	type request struct {
		Name        string                    `json:"name" binding:"required"`
		Description string                    `json:"description" binding:"required"`
		Start       string                    `json:"start" binding:"required"`
		End         string                    `json:"end" binding:"required"`
		Owner       string                    `json:"owner" binding:"required"`
		SharedWith  []int                     `json:"sharedWith" binding:"required"`
		Itinerary   []domain.ItineraryElement `json:"itinerary" binding:"required"`
	}

	type response struct {
		Data domain.Trip `json:"data"`
	}

	return func(c *gin.Context) {
		var newRequest request

		if err := c.ShouldBindJSON(&newRequest); err != nil {
			c.JSON(400, web.NewError(400, "Invalid request"))
			return
		}
		createdTrip, storeErr := t.tripService.Store(c,
			newRequest.Name,
			newRequest.Description,
			newRequest.Start,
			newRequest.End,
			newRequest.Owner,
			newRequest.SharedWith,
			newRequest.Itinerary,
		)

		if storeErr != nil {
			c.JSON(409, web.NewError(409, storeErr.Error()))
			return
		}

		newResponse := response{
			Data: createdTrip,
		}

		c.JSON(201, newResponse)
	}
}

func (t *Trip) Update() gin.HandlerFunc {
	type request struct {
		Name        string                    `json:"name"`
		Description string                    `json:"description"`
		Start       string                    `json:"start"`
		End         string                    `json:"end"`
		Owner       string                    `json:"owner"`
		SharedWith  []int                     `json:"sharedWith"`
		Itinerary   []domain.ItineraryElement `json:"itinerary"`
	}

	type response struct {
		Data domain.Trip `json:"data"`
	}

	return func(c *gin.Context) {
		id := c.Param("id")

		var updReq request

		if err := c.ShouldBindJSON(&updReq); err != nil {
			c.JSON(400, web.NewError(400, err.Error()))
			return
		}

		wUpdated, err := t.tripService.Update(c, id, updReq.Name, updReq.Description, updReq.Start, updReq.End, updReq.Owner, updReq.SharedWith, updReq.Itinerary)
		if err != nil {
			status, _ := strconv.Atoi(err.Error()[0:3])
			c.JSON(status, web.NewError(status, err.Error()))
			return
		}

		res := response{
			Data: wUpdated,
		}
		c.JSON(200, res)
	}
}

func (t *Trip) Delete() gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Param("id")
		delErr := t.tripService.Delete(c, id)

		if delErr != nil {
			c.JSON(404, web.NewError(404, delErr.Error()))
			return
		}

		c.JSON(204, "Trip successfully deleted")
	}
}
