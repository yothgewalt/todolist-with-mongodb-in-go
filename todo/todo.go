package todo

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type Todo struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Title string             `bson:"title,omitempty" json:"title" binding:"required"`
}

type Client struct {
	client *mongo.Client
}

func NewTodo(client *mongo.Client) *Client {
	return &Client{client: client}
}

func (t *Client) newCollection(collectionName string) *mongo.Collection {
	collection := t.client.Database("sample").Collection(collectionName)
	return collection
}

func (t *Client) NewTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		todoListsCollection := t.newCollection("todo_lists")
		todoDocument := bson.M{"title": todo.Title}
		result, err := todoListsCollection.InsertOne(ctx, todoDocument)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		insertedResult := fmt.Sprintf("successfully, the todo has been added (_id: %v)", result.InsertedID)
		c.JSON(http.StatusCreated, gin.H{
			"message": insertedResult,
		})
	}
}

func (t *Client) ReadTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		todoListsCollection := t.newCollection("todo_lists")
		cursor, err := todoListsCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
		}(cursor, ctx)

		var receiver []Todo
		if err = cursor.All(ctx, &receiver); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		var results []Todo
		for _, field := range receiver {
			payload := Todo{
				ID:    field.ID,
				Title: field.Title,
			}

			results = append(results, payload)
		}

		c.JSON(http.StatusFound, results)
	}
}

func (t *Client) DeleteTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id := c.Param("id")
		objId, _ := primitive.ObjectIDFromHex(id)
		todoListsCollection := t.newCollection("todo_lists")

		result, err := todoListsCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		deletedResult := fmt.Sprintf("successfully, the todo has been deleted (count: %v)", result.DeletedCount)
		c.JSON(http.StatusOK, gin.H{
			"message": deletedResult,
		})
	}
}
