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

	userRepository := trip.NewRepository(userCollection)
	userService := trip.NewService(userRepository)
	userHandler := handler.NewTrip(userService)
	userRoutes := router.Group("/api/v1/users")
	{
		userRoutes.GET("/:user_id", userHandler.Get())
		userRoutes.POST("/create_user", userHandler.Store())
		userRoutes.POST("/:user_id/reset_password", userHandler.Update())
		tripRoutes.PATCH("/:user_id", tripHandler.Update())

	}

	router.Run()
}
