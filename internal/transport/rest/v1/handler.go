package restv1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
)

type RESTHandler struct {
	userService userservice.UserSrv
}

func NewRESTHandler(userService userservice.UserSrv) *RESTHandler {
	return &RESTHandler{
		userService: userService,
	}
}

func (h *RESTHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)

		return
	}

	id, err := h.userService.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)

		return
	}
}

func (h *RESTHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)

		return
	}

	user, err := h.userService.GetUser(r.Context(), int64(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)

		return
	}
}

func (h *RESTHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)

		return
	}

	users, err := h.userService.GetUsers(r.Context(), rune(limit))
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)

		return
	}
}

func (h *RESTHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)

		return
	}

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)

		return
	}

	user.ID = int64(id)
	if err := h.userService.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RESTHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)

		return
	}

	if err := h.userService.DeleteUser(r.Context(), int64(id)); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
