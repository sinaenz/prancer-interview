package controller

import (
	"encoding/json"
	"net/http"
	"scratch/internal/service"
)

type HttpController struct {
	CoordinatorService service.Center
}

func (h *HttpController) Move() http.HandlerFunc {
	type in struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		var controllerReq in
		if err := json.NewDecoder(request.Body).Decode(&controllerReq); err != nil {
			panic(err)
		}
		h.CoordinatorService.MoveAgent(request.Context(), &service.Target{
			X: controllerReq.X,
			Y: controllerReq.Y,
		})
	}

}
