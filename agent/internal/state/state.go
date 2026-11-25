package state

import (
	"context"
	"errors"
)

type CursorStore interface {
	Load(ctx context.Context) (uint64, error)
	Save(ctx context.Context, recordID uint64) error
}

var ErrInvalidRecordID = errors.New("record ID must be greater than zero")
