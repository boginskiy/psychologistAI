package msgs

import "fmt"

var (
	MessNeedVerifyAccount             string = "account needs to be verified, please check your email"
	MessExceedingVerificationAttempts string = "number of attempts to verify the account has been exceeded"
	MessOkVerification                string = "verification successful"
	MessOkLogin                       string = "successful login"
	MessOkLogout                      string = "successful logout"

	// Chat
	MessEmptyContent string = "Вы отправили пустое сообщение, повторите свою мысль."
)

var FuncNeedRegistration = func(email string, minutes int) string {
	return fmt.Sprintf("go to '%s' and verify the account for %v minutes", email, minutes)
}

var MapContentMessages = map[string]string{
	"": MessEmptyContent,
}
