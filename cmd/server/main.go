package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gabriel-ballesteros/voyagr-api/cmd/server/handler"
	trip "github.com/gabriel-ballesteros/voyagr-api/internal/trip"
	user "github.com/gabriel-ballesteros/voyagr-api/internal/user"
)

func main() {

	// Set client options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI("mongodb+srv://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASSWORD") + "@" + os.Getenv("MONGO_URL") + "?retryWrites=true&w=majority&appName=voyagr").SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	tripCollection := client.Database("voyagr").Collection("trips")
	userCollection := client.Database("voyagr").Collection("users")

	router := gin.Default()

	tripRepository := trip.NewRepository(tripCollection)
	tripService := trip.NewService(tripRepository)
	tripHandler := handler.NewTrip(tripService)
	tripRoutes := router.Group("/api/v1/trips")
	{
		tripRoutes.GET("", tripHandler.GetAll())
		tripRoutes.GET("/:id", tripHandler.Get())
		tripRoutes.POST("/", tripHandler.Store())
		tripRoutes.PATCH("/:id", tripHandler.Update())
		tripRoutes.DELETE("/:id", tripHandler.Delete())
	}

	userRepository := user.NewRepository(userCollection)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUser(userService)
	userRoutes := router.Group("/api/v1/users")
	{
		userRoutes.GET("/:email", userHandler.Get())
		userRoutes.POST("/create_user", userHandler.Store())
		userRoutes.POST("/:email/reset_password", userHandler.ResetPassword())
		userRoutes.POST("/:email/change_password", userHandler.ChangePassword())
		userRoutes.PATCH("/:email", userHandler.Update())
		userRoutes.DELETE("/:email", userHandler.Delete())

	}

	router.Run()
}
