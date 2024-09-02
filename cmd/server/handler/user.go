package handler

import (
	"strconv"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	user "github.com/gabriel-ballesteros/voyagr-api/internal/user"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
	"github.com/gin-gonic/gin"
)

type User struct {
	userService user.Service
}

func NewUser(u user.Service) *User {
	return &User{
		userService: u,
	}
}

func (u *User) Get() gin.HandlerFunc {
	type response struct {
		Data domain.User `json:"data"`
	}

	return func(c *gin.Context) {
		email := c.Param("email")

		user, werr := u.userService.Get(c, email)
		if werr != nil {
			c.JSON(404, web.NewError(404, "User with email "+email+" not found"))
			return
		}

		res := response{
			Data: user,
		}

		c.JSON(200, res)
		return
	}
}

func (u *User) Store() gin.HandlerFunc {
	type request struct {
		Email string `json:"email" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

	type response struct {
		Data domain.User `json:"data"`
	}

	return func(c *gin.Context) {
		var newRequest request

		if err := c.ShouldBindJSON(&newRequest); err != nil {
			c.JSON(400, web.NewError(400, "Invalid request"))
			return
		}
		createdUser, storeErr := u.userService.Store(c,
			newRequest.Email,
			newRequest.Name)

		if storeErr != nil {
			c.JSON(409, web.NewError(409, storeErr.Error()))
			return
		}

		newResponse := response{
			Data: createdUser,
		}

		c.JSON(201, newResponse)
	}
}

func (u *User) Update() gin.HandlerFunc {
	type request struct {
		Name string `json:"name" binding:"required"`
	}

	type response struct {
		Data domain.User `json:"data"`
	}

	return func(c *gin.Context) {
		email := c.Param("email")

		var updReq request

		if err := c.ShouldBindJSON(&updReq); err != nil {
			c.JSON(400, web.NewError(400, err.Error()))
			return
		}

		uUpdated, err := u.userService.Update(c, email, updReq.Name)
		if err != nil {
			status, _ := strconv.Atoi(err.Error()[0:3])
			c.JSON(status, web.NewError(status, err.Error()))
			return
		}

		res := response{
			Data: uUpdated,
		}
		c.JSON(200, res)
	}
}

func (u *User) ResetPassword() gin.HandlerFunc {

	return func(c *gin.Context) {
		email := c.Param("email")
		err := u.userService.ResetPassword(c, email)
		if err != nil {
			status, _ := strconv.Atoi(err.Error()[0:3])
			c.JSON(status, web.NewError(status, err.Error()))
			return
		}
		c.JSON(200, "Password reseted successfully")
	}
}

func (u *User) ChangePassword() gin.HandlerFunc {
	type request struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	return func(c *gin.Context) {
		var updReq request

		if err := c.ShouldBindJSON(&updReq); err != nil {
			c.JSON(400, web.NewError(400, err.Error()))
			return
		}
		email := c.Param("email")
		err := u.userService.ChangePassword(c, email, updReq.OldPassword, updReq.NewPassword)
		if err != nil {
			status, _ := strconv.Atoi(err.Error()[0:3])
			c.JSON(status, web.NewError(status, err.Error()))
			return
		}
		c.JSON(200, "Password updated successfully")
	}
}

func (u *User) Delete() gin.HandlerFunc {

	return func(c *gin.Context) {
		email := c.Param("email")
		delErr := u.userService.Delete(c, email)

		if delErr != nil {
			c.JSON(404, web.NewError(404, delErr.Error()))
			return
		}

		c.JSON(204, "User deleted successfully")
	}
}
