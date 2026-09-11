package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadAllRequestBody(r *http.Request, item any) (any, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	if err := json.Unmarshal(body, &item); err != nil {
		return nil, fmt.Errorf("failed to deserialization request body: %w", err)
	}
	return item, nil
}
