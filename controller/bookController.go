package controller

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreatBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if book.Title == "" || book.Author == "" || book.Price == 0 {
		http.Error(w, "all field require", http.StatusBadRequest)
	}
	result := db.DB.Create(&book)
	if result.Error != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Failed to encode book", http.StatusInternalServerError)
	}
}
func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Fetching book with ID: %s\n", id)
	var book models.Book
	result := db.DB.First(&book, id)
	if result.Error != nil {
		fmt.Printf("Error retrieving book: %v\n", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve book", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		fmt.Printf("Error encoding book: %v\n", err)
		http.Error(w, "Failed to encode book", http.StatusInternalServerError)
	}
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var book models.Book
	result := db.DB.First(&book, id)
	if result.Error != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := db.DB.Save(&book).Error; err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "Failed to encode book", http.StatusInternalServerError)
	}
}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" {
		http.Error(w, "the id field is require to found the book", http.StatusBadRequest)
		return
	}
	var book models.Book
	result := db.DB.Delete(&book, params["id"])
	if result.Error != nil {
		http.Error(w, "the book id doesn't existe", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("User was deleted")
}
