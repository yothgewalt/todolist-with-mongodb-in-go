package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yongyuth-chuankhuntod/todolists-with-mongodb-in-go/todo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://mongdb:m8Wkv3sZWO8ojUJ4@einer-cluster.aahjt.mongodb.net/myFirstDatabase?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	handler := todo.NewTodo(client)
	r.POST("/todo", handler.NewTask())
	r.GET("/todo", handler.ReadTasks())
	r.DELETE("/todo/:id", handler.DeleteTask())

	err = r.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
