package handlers

import (
	"github.com/dyammarcano/sorteador/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) SponsorHandler(w http.ResponseWriter, r *http.Request) {
	uid := utils.GenerateUUID()
	// Guardar UUID en la base de datos con un timestamp
	// Retornar el UUID
	w.Write([]byte(uid))
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uuid")
	if uid == "" {
		http.Error(w, "UUID no válido", http.StatusBadRequest)
		return
	}

	// Validar el UUID y su expiración en la base de datos
	// Permitir el registro del premio
	w.Write([]byte("Registro exitoso"))
}
