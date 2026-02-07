package utils

import (
	"strings"
)

// ValidateTitle проверяет заголовок чата
func ValidateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return "", &ValidationError{Field: "title", Message: "Title cannot be empty"}
	}

	if len(title) > 200 {
		return "", &ValidationError{Field: "title", Message: "Title must be 200 characters or less"}
	}

	return title, nil
}

// ValidateText проверяет текст сообщения
func ValidateText(text string) (string, error) {
	text = strings.TrimSpace(text)

	if text == "" {
		return "", &ValidationError{Field: "text", Message: "Text cannot be empty"}
	}

	if len(text) > 5000 {
		return "", &ValidationError{Field: "text", Message: "Text must be 5000 characters or less"}
	}

	return text, nil
}

// ValidationError ошибка валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
