package handler

import (
	"api/pkg/response"
	"encoding/json"
	"net/http"

	"api/internal/model"
	"api/internal/service"
)

//recv http

//parse json

// dial service
// return
type CalculatorHandler struct {
	svc *service.CalculatorService
}

func NewCalculatorHandler(s *service.CalculatorService) *CalculatorHandler {
	return &CalculatorHandler{
		svc: s,
	}
}

func (h *CalculatorHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	//recv http request
	//parse json
	//dial service
	//return json
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req model.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	res, err := h.svc.Calculate(req.Left, req.Right, req.Operator)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, model.CalculateResponse{
		Result: res,
	})

}
