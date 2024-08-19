package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"tucows-grill-client/internal/api"
	"tucows-grill-client/internal/models"
)

type Handler struct {
	Client *api.Client
}

func NewHandler(client *api.Client) *Handler {
	return &Handler{Client: client}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var jwt string
	var err error
	if jwt, err = h.Client.Login(creds.Username, creds.Password); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jwt)
}

func (h *Handler) GetIngredientByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ingredientID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	ingredient, err := h.Client.GetIngredientByID(ingredientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredient)
}

func (h *Handler) PostIngredientHandler(w http.ResponseWriter, r *http.Request) {
	var ingredient models.Ingredient
	if err := json.NewDecoder(r.Body).Decode(&ingredient); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Client.PostIngredient(ingredient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Ingredient added successfully!"))
}

// GetTotalCostHandler handles requests to get the total cost for an item_id
func (h *Handler) GetTotalCostHandler(w http.ResponseWriter, r *http.Request) {
    itemIDStr := r.URL.Query().Get("item_id")
    asyncBool := r.URL.Query().Get("async")
    if itemIDStr == "" {
        http.Error(w, "item_id query parameter is required", http.StatusBadRequest)
        return
    }

    itemID, err := strconv.Atoi(itemIDStr)
    if err != nil {
        http.Error(w, "Invalid item_id", http.StatusBadRequest)
        return
    }

    totalCost, err := h.Client.GetTotalCostForItem(itemID, asyncBool)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]float64{"total_cost": totalCost})
}