package restv1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
)

const badRequestJSONFormat = "Invalid JSON format"

type RESTHandler struct {
	service  userservice.UserSrv
	validate *validator.Validate
}

func NewRESTHandler(service userservice.UserSrv, validate *validator.Validate) *RESTHandler {
	return &RESTHandler{
		service:  service,
		validate: validate,
	}
}

func (h *RESTHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.handleBadRequestError(r.Context(), w, badRequestJSONFormat)
		zerolog.Ctx(r.Context()).Err(err).Msgf("validation error: %v", err)

		return
	}

	if err := h.validate.Struct(user); err != nil {
		h.handleBadRequestError(r.Context(), w, err.Error())
		zerolog.Ctx(r.Context()).Err(err).Msgf("validation error: %v", err)

		return
	}

	id, err := h.service.CreateUser(r.Context(), user)
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

	user, err := h.service.GetUser(r.Context(), int64(id))
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

	users, err := h.service.GetUsers(r.Context(), rune(limit))
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
		h.handleBadRequestError(r.Context(), w, badRequestJSONFormat)

		return
	}

	if err := h.validate.Struct(user); err != nil {
		h.handleBadRequestError(r.Context(), w, err.Error())
		zerolog.Ctx(r.Context()).Err(err).Msgf("validation error: %v", err)

		return
	}

	user.ID = int64(id)
	if err := h.service.UpdateUser(r.Context(), user); err != nil {
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

	if err := h.service.DeleteUser(r.Context(), int64(id)); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RESTHandler) handleBadRequestError(ctx context.Context, res http.ResponseWriter, errorMessage string) {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusBadRequest)

	errResponse := entity.ErrResponse{Error: errorMessage}
	if err := json.NewEncoder(res).Encode(errResponse); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("encode reserve period or date")
	}
}
