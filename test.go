package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/boginskiy/psychologistAI/internal/web/renders"
)

const (
	HostEmail = "smtp.gmail.com"
	PortEmail = 587
	FromEmail = "gophkeeper@gmail.com"
	Password  = "upiplnvviujgnevc"

	Subject = "Подтверждение email"
	Retry   = 3
	TimeOut = 200 * time.Millisecond
)

// func main() {

// 	for range 0 {
// 		fmt.Println("ddd")
// 	}

// 	// _ = db
// }

func (h *ChatHandler) ShowChat(w http.ResponseWriter, r *http.Request) {
	// 1. Достаём пользователя из сессии/контекста
	infoUser, err := h.getInfoUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 2. Получаем чат (или пустой)
	userChat, err := h.chatService.GetOrCreate(r.Context(), infoUser.UserID)
	if err != nil {
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	// 3. Готовим историю для шаблона
	history := make([]map[string]any, 0, len(userChat.Messages))
	for _, msg := range userChat.Messages {
		history = append(history, map[string]any{
			"Role": msg.Role, // "user" или "assistant"
			"Text": msg.Text, // содержимое
			"Time": formatTime(msg.CreatedAt),
		})
	}

	// 4. Рендерим
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/chat.html",
		"templates/partials/message_user.html",
		"templates/partials/message_ai.html",
	)
	if err != nil {
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	data := map[string]any{
		"History": history,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("chat template error: %v", err)
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04")
}
