package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"todo-server/helpers"
	"todo-server/models"
)

type UserController struct{}

type ResponseOutput struct {
	User  models.User
	Token string
}

func (u UserController) SignupUser(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		User := models.User{}
		UserCred := models.UserCred{}
		error := helpers.CustomError{}
		userCollection := db.Database("todoapp").Collection("user")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		json.NewDecoder(r.Body).Decode(&UserCred)

		result, insertErr := userCollection.InsertOne(ctx, UserCred)
		if insertErr != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Add new User in database! \n"+insertErr.Error())
			return
		}

		User.FirstName = UserCred.FirstName
		User.LastName = UserCred.LastName
		User.Email = UserCred.Email
		User.ID = result.InsertedID.(primitive.ObjectID)

		payload := helpers.Payload{
			FirstName: User.FirstName,
			LastName:  User.LastName,
			Email:     User.Email,
		}
		token, err := helpers.GenerateJwtToken(payload)
		if err != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Generate New JWT Token!")
			return
		}

		helpers.RespondWithJSON(w, ResponseOutput{
			Token: token,
			User:  User,
		})
	}
}

func (u UserController) LoginUser(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		User := models.User{}
		error := helpers.CustomError{}
		userCollection := db.Database("todoapp").Collection("user")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		type Credentials struct {
			Email    string
			Password string
		}
		credentials := Credentials{}
		json.NewDecoder(r.Body).Decode(&credentials)

		results := userCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&User)
		if results != nil {
			error.ApiError(w, http.StatusNotFound, "Invalid Username/Email, Please Signup!")
			return
		}

		if User.Password != credentials.Password {
			error.ApiError(w, http.StatusNotFound, "Invalid Credentials!")
			return
		}

		payload := helpers.Payload{
			FirstName: User.FirstName,
			LastName:  User.LastName,
			Email:     User.Email,
		}

		token, err := helpers.GenerateJwtToken(payload)
		if err != nil {
			error.ApiError(w, http.StatusInternalServerError, "Failed To Generate New JWT Token!")
			return
		}

		User.Password = ""
		helpers.RespondWithJSON(w, ResponseOutput{
			Token: token,
			User:  User,
		})
	}
}
