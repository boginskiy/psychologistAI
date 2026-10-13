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




	if err != nil {
		log.Printf("error: %v\n", err) // + logger

		body := &response.InfoBody{}

		switch {
		// Credentials
		case errors.Is(err, errs.ErrSession):
			body.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)

		case errors.Is(err, errs.ErrUsingToken), errors.Is(err, errs.ErrLegitimacyUser):
			body.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)

		case errors.Is(err, errs.ErrTimeLiveSession), errors.Is(err, errs.ErrCompareToken):
			body.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)

		// Server
		case errors.Is(err, errs.ErrServer), errors.Is(err, errs.ErrUpdateDBAfterRefresh):
			body.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)

		default:
			body = response.NewInfoBodyWithErr(err, http.StatusBadRequest)
		}
		h.Responder.SendResponse(w, body)
		return
	}


		// Errors
	if err != nil {
		log.Printf("error: %v\n", err) // + logger

		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			renders.RenderError(w, "401", http.StatusUnauthorized, nil)

		// Verification
		case errors.Is(err, errs.ErrVerification):
			renders.RenderError(w, "403", http.StatusForbidden, nil)
		case errors.Is(err, errs.ErrAttemptsVerification):
			renders.RenderError(w, "429", http.StatusTooManyRequests, nil)

		// Server
		case errors.Is(err, errs.ErrServer):
			renders.RenderError(w, "500", http.StatusInternalServerError, nil)

		default:
			renders.RenderError(w, "400", http.StatusBadRequest, nil)
		}

		return
	}