package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dyammarcano/sorteador/internal/embed"
	"github.com/dyammarcano/sorteador/internal/models"
	"github.com/dyammarcano/sorteador/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"html/template"
	"math/rand"
	"net/http"
	"time"
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

	// Redirigir a la página de registro del sponsor
	http.Redirect(w, r, fmt.Sprintf("/sponsor/register/%s", uid), http.StatusSeeOther)
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

func (h *Handler) SponsorSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.FormValue("uuid")
	lid := r.FormValue("ulid")
	name := r.FormValue("name")
	prize := r.FormValue("prize")

	// Guardar los datos en la base de datos
	sponsor := models.Sponsor{
		UUID:      uid,
		ULID:      lid,
		Name:      name,
		Prize:     prize,
		Timestamp: time.Now().Unix(),
	}

	// query database to find ulid if exists and return error
	var u string
	if err := h.db.Get(&u, "SELECT ulid FROM sponsors WHERE ulid = $1", sponsor.ULID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.ErrorJSON(w, fmt.Sprintf("error fetching sponsor: %v", err))
			return
		}
	}

	if u != "" {
		utils.ErrorJSON(w, fmt.Sprintf("sponsor operation already exists"))
		return
	}

	if _, err := h.db.Exec("INSERT INTO sponsors (uuid, ulid, name, prize, timestamp) VALUES ($1, $2, $3, $4, $5)", sponsor.UUID, sponsor.ULID, sponsor.Name, sponsor.Prize, sponsor.Timestamp); err != nil {
		utils.ErrorJSON(w, fmt.Sprintf("error inserting sponsor: %v", err))
	}

	// Aquí puedes agregar lógica para guardar el sponsor en la base de datos

	w.Write([]byte(fmt.Sprintf("Sponsor %s registered with prize %s", name, prize)))
}

func (h *Handler) SponsorRegisterHandler(w http.ResponseWriter, r *http.Request) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	uid := chi.URLParam(r, "uuid")
	static := embed.GetStatic()
	tmpl, err := template.ParseFS(static, "templates/sponsor_register.html")
	if err != nil {
		utils.ErrorJSON(w, fmt.Sprintf("error parsing template: %v", err))
		return
	}

	if err = tmpl.Execute(w, map[string]string{"UUID": uid, "ULID": id}); err != nil {
		utils.ErrorJSON(w, fmt.Sprintf("error rendering template: %v", err))
		return
	}
}

func (h *Handler) StaticHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(embed.GetStatic())).ServeHTTP(w, r)
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	static := embed.GetStatic()
	tmpl, err := template.ParseFS(static, "templates/index.html")
	if err != nil {
		utils.ErrorJSON(w, fmt.Sprintf("error parsing template: %v", err))
		return
	}

	if err = tmpl.Execute(w, map[string]string{"ULID": id, "Title": "My Sorteador"}); err != nil {
		utils.ErrorJSON(w, fmt.Sprintf("error rendering template: %v", err))
		return
	}
}
