package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"todo-server/helpers"
	"todo-server/models"
)

type TodoController struct{}

var error = helpers.CustomError{}

func (t TodoController) GetTodoListByUserId(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		TodoList := []models.Todo{}
		todoCollection := db.Database("todoapp").Collection("todo")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["user"])

		cur, getErrors := todoCollection.Find(ctx, bson.M{"user": id})
		if getErrors != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Fetch ToDo Tasks from database!")
			return
		}
		for cur.Next(context.Background()) {
			result := models.Todo{}
			e := cur.Decode(&result)
			if e != nil {
				log.Fatal(e)
			}
			TodoList = append(TodoList, result)
		}
		if err := cur.Err(); err != nil {

			log.Fatal(err)
		}
		helpers.RespondWithJSON(w, TodoList)
	}
}

func (t TodoController) AddNewTodo(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		TodoTask := models.Todo{}
		json.NewDecoder(r.Body).Decode(&TodoTask)
		todoCollection := db.Database("todoapp").Collection("todo")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, getErrors := todoCollection.InsertOne(ctx, TodoTask)
		if getErrors != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Add new Task in database!")
			return
		}

		TodoTask.ID = result.InsertedID.(primitive.ObjectID)
		helpers.RespondWithJSON(w, TodoTask)
	}
}

func (t TodoController) UpdateTodo(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		TodoTask := models.Todo{}
		params := mux.Vars(r)
		json.NewDecoder(r.Body).Decode(&TodoTask)
		todoCollection := db.Database("todoapp").Collection("todo")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, _ := primitive.ObjectIDFromHex(params["id"])
		update := bson.M{"taskname": TodoTask.Taskname, "task": TodoTask.Task, "targetdate": TodoTask.Targetdate, "status": TodoTask.Status}

		_, getErrors := todoCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
		if getErrors != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Add new Task in database!")
			return
		}

		//TodoTask.ID = result.UpsertedID.(primitive.ObjectID)
		helpers.RespondWithJSON(w, TodoTask)
	}
}

func (t TodoController) TaskComplete(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		TodoTask := models.Todo{}
		params := mux.Vars(r)
		todoCollection := db.Database("todoapp").Collection("todo")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, _ := primitive.ObjectIDFromHex(params["id"])
		_, getErrors := todoCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": true}})
		if getErrors != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed to Update Task from the database!")
			return
		}
		helpers.RespondWithJSON(w, TodoTask)
	}
}

func (t TodoController) UndoTask(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		TodoTask := models.Todo{}
		params := mux.Vars(r)
		todoCollection := db.Database("todoapp").Collection("todo")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, _ := primitive.ObjectIDFromHex(params["id"])
		_, getErrors := todoCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": false}})
		if getErrors != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed to UnDo Task from the database!")
			return
		}
		helpers.RespondWithJSON(w, TodoTask)
	}
}

func (t TodoController) DeleteTodo(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		TodoTask := models.Todo{}
		params := mux.Vars(r)
		todoCollection := db.Database("todoapp").Collection("todo")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, _ := primitive.ObjectIDFromHex(params["id"])
		_, getErrors := todoCollection.DeleteOne(ctx, bson.M{"_id": id})
		if getErrors != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed to Delete Item from the database!")
			return
		}
		helpers.RespondWithJSON(w, TodoTask)
	}
}
