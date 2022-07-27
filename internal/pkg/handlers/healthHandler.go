package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// HealthCheckHandler godoc
// @Description Check if server is healthy
// @Tags health
// @Produce application/json
// @Success 200 {object} CommonResponse "'alive'"
// @Router /health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var commonResponse CommonResponse
	commonResponse.Message = "alive"
	err := json.NewEncoder(w).Encode(commonResponse)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
