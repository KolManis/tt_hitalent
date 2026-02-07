package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/KolManis/tt_hitalent/internal/models"
	"github.com/KolManis/tt_hitalent/pkg/utils"
	"gorm.io/gorm"
)

type ChatHandler struct {
	db *gorm.DB
}

func NewChatHandler(db *gorm.DB) *ChatHandler {
	return &ChatHandler{db: db}
}

// CreateChat - POST /chats
func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validatedTitle, err := utils.ValidateTitle(request.Title)
	if err != nil {
		writeValidationError(w, err)
		return
	}

	chat := models.Chat{
		Title:     validatedTitle,
		CreatedAt: time.Now(),
	}

	if err := h.db.Create(&chat).Error; err != nil {
		writeError(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

// GetChat - GET /chats/{id}
func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		writeError(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Получаем параметр limit из query string
	limit := 20 // значение по умолчанию
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			if l > 0 && l <= 100 {
				limit = l
			}
		}
	}

	var chat models.Chat
	if err := h.db.First(&chat, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			writeError(w, "Chat not found", http.StatusNotFound)
		} else {
			writeError(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	var messages []models.Message
	h.db.Where("chat_id = ?", id).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages)

	// Разворачиваем порядок сообщений (был DESC, нужен ASC)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	response := struct {
		ID        uint             `json:"id"`
		Title     string           `json:"title"`
		CreatedAt time.Time        `json:"created_at"`
		Messages  []models.Message `json:"messages"`
	}{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
		Messages:  messages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteChat - DELETE /chats/{id}
func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		writeError(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Каскадное удаление (сообщения удалятся автоматически благодаря FOREIGN KEY ON DELETE CASCADE)
	result := h.db.Delete(&models.Chat{}, id)
	if result.Error != nil {
		writeError(w, "Failed to delete chat", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		writeError(w, "Chat not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func writeValidationError(w http.ResponseWriter, err error) {
	if valErr, ok := err.(*utils.ValidationError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": valErr.Error(),
			"field": valErr.Field,
		})
	} else {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
}
