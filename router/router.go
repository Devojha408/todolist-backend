package router

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	controller "todo-server/controllers"
	middleware "todo-server/middlewares"
)

func RegisterRoutes(db *mongo.Client) *mux.Router {
	router := mux.NewRouter()

	TodoController := controller.TodoController{}

	router.HandleFunc("/todo/all/{user}", middleware.CheckAuth(TodoController.GetTodoListByUserId(db))).Methods("GET")
	router.HandleFunc("/todo", middleware.CheckAuth(TodoController.AddNewTodo(db))).Methods("POST")
	router.HandleFunc("/todo/update/{id}", middleware.CheckAuth(TodoController.UpdateTodo(db))).Methods("PUT")
	router.HandleFunc("/todo/done/{id}", middleware.CheckAuth(TodoController.TaskComplete(db))).Methods("PUT")
	router.HandleFunc("/todo/undo/{id}", middleware.CheckAuth(TodoController.UndoTask(db))).Methods("PUT")
	router.HandleFunc("/todo/{id}", middleware.CheckAuth(TodoController.DeleteTodo(db))).Methods("DELETE")

	UserController := controller.UserController{}

	router.HandleFunc("/auth/login", UserController.LoginUser(db)).Methods("POST")
	router.HandleFunc("/auth/signup", UserController.SignupUser(db)).Methods("POST")

	return router
}
