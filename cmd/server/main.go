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

	clientOptions := options.Client().ApplyURI("mongodb+srv://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASSWORD") + "@" + os.Getenv("MONGO_URL") + "?retryWrites=true&w=majority&appName=voyagr")

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
	collection := client.Database("voyagr").Collection("trips")

	router := gin.Default()

	// WAREHOUSES
	repository := trip.NewRepository(collection)
	service := trip.NewService(repository)
	handler := handler.NewTrip(service)
	routes := router.Group("/api/v1/trips")
	{
		routes.GET("/", handler.GetAll())
		routes.GET("/:id", handler.Get())
		routes.POST("/", handler.Store())
		routes.PATCH("/:id", handler.Update())
		routes.DELETE("/:id", handler.Delete())
	}

	router.Run()
}
