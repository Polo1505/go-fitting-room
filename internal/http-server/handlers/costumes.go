package handlers

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/Polo1505/go-fitting-room/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type handlers struct {
	storage *postgresql.Storage
	logger  *slog.Logger
}

// New создаёт новый экземпляр handlers
func New(storage *postgresql.Storage, logger *slog.Logger) *handlers {
	return &handlers{
		storage: storage,
		logger:  logger,
	}
}

// CreateCostume обрабатывает POST /costumes
func (h *handlers) CreateCostume(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CreateCostume"

	var costume postgresql.Costume
	if err := json.NewDecoder(r.Body).Decode(&costume); err != nil {
		h.logger.Error("Failed to decode request body", slog.String("op", op), slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.storage.CreateCostume(&costume); err != nil {
		h.logger.Error("Failed to create costume", slog.String("op", op), slog.String("error", err.Error()))
		http.Error(w, "Failed to create costume", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(costume); err != nil {
		h.logger.Error("Failed to encode response", slog.String("op", op), slog.String("error", err.Error()))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetCostume обрабатывает GET /costumes/{id}
func (h *handlers) GetCostume(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetCostume"

	idStr := chi.URLParam(r, "id")
	h.logger.Debug("Received ID", slog.String("op", op), slog.String("id", idStr)) // Добавлено для отладки
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid costume ID", slog.String("op", op), slog.String("id", idStr), slog.String("error", err.Error()))
		http.Error(w, "Invalid costume ID", http.StatusBadRequest)
		return
	}

	costume, err := h.storage.GetCostume(id)
	if err != nil {
		h.logger.Error("Failed to get costume", slog.String("op", op), slog.String("error", err.Error()))
		if err.Error() == "costume not found" {
			http.Error(w, "Costume not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get costume", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(costume); err != nil {
		h.logger.Error("Failed to encode response", slog.String("op", op), slog.String("error", err.Error()))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetAllCostumes обрабатывает GET /costumes
func (h *handlers) GetAllCostumes(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetAllCostumes"

	costumes, err := h.storage.GetAllCostumes()
	if err != nil {
		h.logger.Error("Failed to get all costumes", slog.String("op", op), slog.String("error", err.Error()))
		http.Error(w, "Failed to get costumes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(costumes); err != nil {
		h.logger.Error("Failed to encode response", slog.String("op", op), slog.String("error", err.Error()))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateCostume обрабатывает PUT /costumes/{id}
func (h *handlers) UpdateCostume(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.UpdateCostume"

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid costume ID", slog.String("op", op), slog.String("id", idStr), slog.String("error", err.Error()))
		http.Error(w, "Invalid costume ID", http.StatusBadRequest)
		return
	}

	var updatedCostume postgresql.Costume
	if err := json.NewDecoder(r.Body).Decode(&updatedCostume); err != nil {
		h.logger.Error("Failed to decode request body", slog.String("op", op), slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.storage.UpdateCostume(id, &updatedCostume); err != nil {
		h.logger.Error("Failed to update costume", slog.String("op", op), slog.String("error", err.Error()))
		if err.Error() == "costume not found" {
			http.Error(w, "Costume not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update costume", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteCostume обрабатывает DELETE /costumes/{id}
func (h *handlers) DeleteCostume(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.DeleteCostume"

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid costume ID", slog.String("op", op), slog.String("id", idStr), slog.String("error", err.Error()))
		http.Error(w, "Invalid costume ID", http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteCostume(id); err != nil {
		h.logger.Error("Failed to delete costume", slog.String("op", op), slog.String("error", err.Error()))
		if err.Error() == "costume not found" {
			http.Error(w, "Costume not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete costume", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
