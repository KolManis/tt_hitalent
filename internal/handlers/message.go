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

type MessageHandler struct {
	db *gorm.DB
}

func NewMessageHandler(db *gorm.DB) *MessageHandler {
	return &MessageHandler{db: db}
}

// CreateMessage - POST /chats/{id}/messages
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	chatIDStr := r.PathValue("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil || chatID <= 0 {
		writeError(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var chat models.Chat
	if err := h.db.First(&chat, chatID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			writeError(w, "Chat not found", http.StatusNotFound)
		} else {
			writeError(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	var req struct {
		Text string `json:"text"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	validatedText, err := utils.ValidateText(req.Text)
	if err != nil {
		writeValidationError(w, err)
		return
	}

	message := models.Message{
		ChatID:    uint(chatID),
		Text:      validatedText,
		CreatedAt: time.Now(),
	}

	if err := h.db.Create(&message).Error; err != nil {
		writeError(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}
