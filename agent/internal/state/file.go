package state

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FileCursorStore struct {
	path string
}

type cursorState struct {
	LastRecordID uint64 `json:"lastRecordId"`
}

func NewFileCursorStore(path string) (*FileCursorStore, error) {
	if path == "" {
		return nil, fmt.Errorf("state file path is required")
	}

	return &FileCursorStore{
		path: path,
	}, nil
}

func (s *FileCursorStore) Load(ctx context.Context) (uint64, error) {
	if ctx == nil {
		return 0, fmt.Errorf("context is required")
	}

	if err := ctx.Err(); err != nil {
		return 0, err
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}

		return 0, fmt.Errorf("read cursor state: %w", err)
	}

	var state cursorState

	if err := json.Unmarshal(data, &state); err != nil {
		return 0, fmt.Errorf("decode cursor state: %w", err)
	}

	return state.LastRecordID, nil
}

func (s *FileCursorStore) Save(
	ctx context.Context,
	recordID uint64,
) error {
	if ctx == nil {
		return fmt.Errorf("context is required")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if recordID == 0 {
		return ErrInvalidRecordID
	}

	data, err := json.Marshal(cursorState{
		LastRecordID: recordID,
	})
	if err != nil {
		return fmt.Errorf("encode cursor state: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}

	tempPath := s.path + ".tmp"

	if err := os.WriteFile(tempPath, data, 0600); err != nil {
		return fmt.Errorf("write cursor state: %w", err)
	}

	if err := os.Rename(tempPath, s.path); err != nil {
		return fmt.Errorf("replace cursor state: %w", err)
	}

	return nil
}

var _ CursorStore = (*FileCursorStore)(nil)
