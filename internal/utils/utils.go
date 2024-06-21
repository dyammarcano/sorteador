package utils

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func IsUUIDExpired(timestamp int64) bool {
	return time.Now().Unix()-timestamp > 86400
}

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ErrorJSON(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	ResponseJSON(w, map[string]string{"error": message})
}
