package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/boginskiy/psychologistAI/internal/api/request"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/service"
	adaptersWEB "github.com/boginskiy/psychologistAI/internal/web/adapters"
	"github.com/boginskiy/psychologistAI/internal/web/renders"
	"github.com/go-chi/chi"
)

type HomeHandler struct {
	AuthService service.AuthService
	ChatService service.ChatService
}

func NewHomeHandler(authService service.AuthService, chatService service.ChatService) *HomeHandler {
	return &HomeHandler{
		AuthService: authService,
		ChatService: chatService,
	}
}

func (h *HomeHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Use(middleware.AuthWebMiddleware(h.AuthService))
	r.Get("/", h.Starter)
	r.Post("/chat", h.Sender)
}

func (h *HomeHandler) Sender(w http.ResponseWriter, r *http.Request) {
	infoUser, err := adaptersWEB.ToInfoUserFromRequest(r)

	// Не смогли достать сообщение от пользователя. Возможна ошибка сервера
	// Нет активных действий. Состояние чата остается текущим.
	if err != nil || infoUser.Message == "" {
		return
	}

	fmt.Println("0000")

	// Service
	userChat, err := h.ChatService.Send(r.Context(), infoUser)

	fmt.Println("11111")

	// Errors
	if err != nil {
		log.Printf("service error: %v\n", err) // + logger

		switch {
		case errors.Is(err, errs.ErrReadChat):
			http.Error(w, "error read chat", http.StatusInternalServerError)

		// Одиночный ответ пользователя
		case errors.Is(err, errs.ErrAIResponse):

			last := userChat.LastN(1)
			if len(last) < 1 {
				http.Error(w, "unexpected messages count", http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles(
				"templates/partials/message_user.html",
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			userMsg := map[string]any{
				"Role": "user",
				"Text": last[0].Text,
				"Time": formatTime(last[0].CreatedAt),
			}

			if err := tmpl.ExecuteTemplate(w, "message_user", userMsg); err != nil {
				log.Printf("render user message: %v", err)
			}

		default:
		}
		return
	}

	fmt.Println("1111")

	// Берём последние два сообщения (user + ai)
	last := userChat.LastN(2)
	if len(last) < 2 {
		http.Error(w, "unexpected messages count", http.StatusInternalServerError)
		return
	}

	fmt.Println("22222")

	// Парсим только партиалы — base.html НЕ нужен
	tmpl, err := template.ParseFiles(
		"templates/partials/message_user.html",
		"templates/partials/message_ai.html",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("3333")

	// 7. Рендерим оба фрагмента в один ответ
	userMsg := map[string]any{
		"Role": "user",
		"Text": last[0].Text,
		"Time": formatTime(last[0].CreatedAt),
	}
	aiMsg := map[string]any{
		"Role": "assistant",
		"Text": last[1].Text,
		"Time": formatTime(last[1].CreatedAt),
	}

	if err := tmpl.ExecuteTemplate(w, "message_user", userMsg); err != nil {
		log.Printf("render user message: %v", err)
	}
	if err := tmpl.ExecuteTemplate(w, "message_ai", aiMsg); err != nil {
		log.Printf("render ai message: %v", err)
	}
}

func (h *HomeHandler) Starter(w http.ResponseWriter, r *http.Request) {
	infoUser, ok := request.GetInfoUserFromContext(r.Context())

	// Получаем историю чата для зарегистрированного пользователя
	// Получаем пустой чат для незарегистрированного пользователя
	userChat, err := h.ChatService.GetOrCreate(r.Context(), infoUser)
	if err != nil {
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	// Готовим историю для шаблона
	history := make([]map[string]any, 0, len(userChat.Messages))
	for _, msg := range userChat.Messages {
		history = append(history, map[string]any{
			"Role": msg.Role, // "user" или "assistant"
			"Text": msg.Text, // содержимое
			"Time": formatTime(msg.CreatedAt),
		})
	}

	// Рендерим
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/index.html",
		"templates/chat.html",
		"templates/chat/message_user.html",
		"templates/chat/message_ai.html",
	)

	if err != nil {
		// + logger
		log.Printf("template error: %v", err)
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	data := map[string]any{
		"History":    history,
		"AuthStatus": h.takeAuthStatus(ok),
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

// takeAuthStatus
func (h *HomeHandler) takeAuthStatus(isUser bool) string {
	status := "LOG IN"
	if isUser {
		status = "LOG OUT"
	}
	return status
}

// formatTime
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04")
}

// TODO
// Нужно хранить запрос/ответ и каждый раз отдавать его клиенту.

// Нужно точно хранилище !

// Процесс будет зависать пока модель не подготовит ответ?
// Или мы отдаем сразу ответ пользователю и как только ответ от модели придет
// мы его отдадим пользователю, но как мы поймем какому пользователю отдавать?
// Хранить ip ?
