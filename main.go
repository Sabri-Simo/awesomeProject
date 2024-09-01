package main

import (
	"awesomeProject/controller"
	"awesomeProject/db"
	"awesomeProject/models"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	db.Connect()

	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.Book{})
	db.DB.AutoMigrate(&models.Card{})

	router := mux.NewRouter()

	router.HandleFunc("/users", controller.CreateUser).Methods("POST")
	router.HandleFunc("/users", controller.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", controller.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/book", controller.CreatBook).Methods("POST")
	router.HandleFunc("/book/{id}", controller.GetBook).Methods("GET")
	router.HandleFunc("/updatebook/{id}", controller.UpdateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", controller.DeleteBook).Methods("DELETE")
	router.HandleFunc("/user/{id}/addbook", controller.AddBook).Methods("POST")
	router.HandleFunc("/user/{id}/deletebook", controller.DeleteBookUser).Methods("DELETE")
	router.HandleFunc("/user/{id}/total", controller.CardTotalPrice).Methods("GET")

	http.ListenAndServe(":8080", router)
}
