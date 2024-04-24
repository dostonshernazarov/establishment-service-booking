package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BaseUseCase struct{}

func (u *BaseUseCase) Error(msg string, err error) error {
	if len(strings.TrimSpace(msg)) != 0 {
		return fmt.Errorf("%v: %w", msg, err)
	}
	return err
}

func (u *BaseUseCase) beforeRequest(id *string, createdAt *time.Time, updatedAt *time.Time, deletedAt *time.Time) {
	if id != nil {
		*id = uuid.New().String() 
	}

	if createdAt != nil {
		*createdAt = time.Now().UTC()
	}

	if updatedAt != nil {
		*updatedAt = time.Now().UTC()
	}

	if deletedAt != nil {
		*deletedAt = time.Now().UTC()
	}
}
