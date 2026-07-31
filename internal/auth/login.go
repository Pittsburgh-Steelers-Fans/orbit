package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// Handler serves registration and login requests backed by an in-memory store.
type Handler struct {
	mu          sync.RWMutex
	secret      string
	ttl         time.Duration
	credentials map[string]credential
}

type credential struct {
	UserID       string
	Email        string
	PasswordHash string
	Roles        []string
}

type registerRequest struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

// NewHandler creates an auth handler using JWT_SECRET when secret is empty.
func NewHandler(secret string, ttl time.Duration) *Handler {
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &Handler{secret: secret, ttl: ttl, credentials: make(map[string]credential)}
}

// Routes mounts auth endpoints on a chi router.
func (h *Handler) Routes(r chi.Router) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

// Register stores a user's hashed password for later login.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.Context().Err(); err != nil {
		h.writeError(w, http.StatusRequestTimeout, err)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	userID := strings.TrimSpace(req.UserID)
	if userID == "" || email == "" || req.Password == "" {
		h.writeError(w, http.StatusBadRequest, errors.New("auth: user_id, email, and password are required"))
		return
	}
	roles := req.Roles
	if len(roles) == 0 {
		roles = []string{"user"}
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.credentials[email]; exists {
		h.writeError(w, http.StatusConflict, errors.New("auth: email already registered"))
		return
	}
	h.credentials[email] = credential{UserID: userID, Email: email, PasswordHash: hashPassword(req.Password), Roles: roles}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}

// Login verifies credentials and writes a signed JWT response.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.Context().Err(); err != nil {
		h.writeError(w, http.StatusRequestTimeout, err)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))

	h.mu.RLock()
	stored, ok := h.credentials[email]
	h.mu.RUnlock()
	if !ok || stored.PasswordHash != hashPassword(req.Password) {
		h.writeError(w, http.StatusUnauthorized, errors.New("auth: invalid credentials"))
		return
	}
	token, err := Issue(h.secret, stored.UserID, h.ttl)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tokenResponse{Token: token})
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func (h *Handler) writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
