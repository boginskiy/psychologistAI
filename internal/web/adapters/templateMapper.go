package adapters

func ToMapRegistrTemplate(email string, liveTimeToken int) map[string]any {
	return map[string]any{
		"Email":         email,
		"LiveTimeToken": liveTimeToken,
	}
}

func ToMapStartTemplate(isUser bool) map[string]any {
	statusAuth := "LOG IN"
	if isUser {
		statusAuth = "LOG OUT"
	}

	return map[string]any{
		"StatusAuth": statusAuth,
	}
}
