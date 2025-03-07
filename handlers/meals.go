package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mjande/meals-microservice/models"
	"github.com/mjande/meals-microservice/utils"
)

type MealResponse struct {
	Message string        `json:"message"`
	Data    []models.Meal `json:"data"`
}

// Handles getting a list of meals for a user.
func GetMeals(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	meals, err := models.ListMealsByDate(r.Context(), params.Get("start"), params.Get("end"))
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := MealResponse{
		Data: meals,
	}

	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Handles creating a new meal.
func PostMeal(w http.ResponseWriter, r *http.Request) {
	var meal models.Meal
	err := json.NewDecoder(r.Body).Decode(&meal)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := models.CreateMeal(r.Context(), meal)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	meal, err = models.FindMeal(r.Context(), id)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := MealResponse{
		Message: "Meal successfully created!",
		Data:    []models.Meal{meal},
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println()
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Handles deleting a meal by id.
func DeleteMeal(w http.ResponseWriter, r *http.Request) {
	// Extract id from request
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = models.DeleteMeal(r.Context(), id)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := MealResponse{
		Message: "Meal successfully deleted",
	}

	w.WriteHeader(http.StatusNoContent)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}
