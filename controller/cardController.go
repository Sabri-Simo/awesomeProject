package controller

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func AddBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var input struct {
		BookID uint `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var card models.Card
	if err := db.DB.Where("user_id = ?", userID).First(&card).Error; err != nil {
		card = models.Card{
			UserID: user.ID,
		}
		if err := db.DB.Create(&card).Error; err != nil {
			http.Error(w, "Failed to create card", http.StatusInternalServerError)
			return
		}
	}

	var book models.Book
	if err := db.DB.First(&book, input.BookID).Error; err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if err := db.DB.Model(&card).Association("Books").Append(&book); err != nil {
		http.Error(w, "Failed to add book to card", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Preload("Books").First(&card).Error; err != nil {
		http.Error(w, "Failed to retrieve updated card", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(card)
}

func DeleteBookUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var input struct {
		BookID uint `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.Preload("Cards").First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var book models.Book
	if err := db.DB.First(&book, input.BookID).Error; err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	for i := range user.Cards {
		db.DB.Model(&user.Cards[i]).Association("Books").Delete(&book)
	}

	db.DB.Preload("Cards.Books").First(&user, userID)
	json.NewEncoder(w).Encode(user)
}
func CardTotalPrice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]
	var user models.User
	db.DB.Preload("Cards.Books").First(&user, userID)
	var totalPrice float32
	for _, card := range user.Cards {
		for _, book := range card.Books {
			totalPrice += book.Price
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Total Price of Books in Cards: %.2f", totalPrice)

}
