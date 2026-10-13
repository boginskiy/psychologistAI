package errs

import "net/http"

type Catcher interface {
	CatchWebErr(w http.ResponseWriter, err error)
}
