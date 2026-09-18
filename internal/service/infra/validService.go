package infra

import (
	"context"
	"fmt"
)

type ValidService struct {
}

func NewValidService(ctx context.Context) *ValidService {
	return &ValidService{}
}

func (s *ValidService) CheckNotEmptyStrField(name, field string) error {
	if field == "" {
		return fmt.Errorf("field '%s' cannot be empty string", name)
	}
	return nil
}
