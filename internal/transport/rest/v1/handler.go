package restv1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog"

	servicev1 "github.com/powerdigital/go-micro/internal/service/v1"
)

const (
	badJSONRequestFormat = "Invalid JSON message"
	resultGettingError   = "Result getting failed"
)

type helloRequestMsg struct {
	Name string `json:"name"`
}

type helloResponseMsg struct {
	HelloName string `json:"name"`
}

type errResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	service servicev1.GreetingService
}

func NewHandler(service servicev1.GreetingService) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) GetHello(res http.ResponseWriter, req *http.Request) {
	jsonRaw, err := io.ReadAll(req.Body)
	if err != nil {
		zerolog.Ctx(req.Context()).Err(err).Msg("get request")

		return
	}

	var requestMsg helloRequestMsg
	if err := json.Unmarshal(jsonRaw, &requestMsg); err != nil {
		h.handleBadRequestError(req.Context(), res, badJSONRequestFormat)

		return
	}

	result, err := h.service.GetHello(requestMsg.Name)
	if err != nil {
		zerolog.Ctx(req.Context()).Err(err).Msg("get hello request")

		h.handleBadRequestError(req.Context(), res, resultGettingError)

		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	responseMsg := helloResponseMsg{
		HelloName: result,
	}

	if err := json.NewEncoder(res).Encode(responseMsg); err != nil {
		zerolog.Ctx(req.Context()).Err(err).Msg("encode reserve period")
	}
}

func (h Handler) handleBadRequestError(ctx context.Context, res http.ResponseWriter, errorMessage string) {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusBadRequest)

	errResponse := errResponse{Error: errorMessage}
	if err := json.NewEncoder(res).Encode(errResponse); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("encode reserve period")
	}
}
