package state

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileCursorStoreLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cursor.json")

	store, err := NewFileCursorStore(path)
	if err != nil {
		t.Fatalf("expected store creation to succeed: %v", err)
	}

	recordID, err := store.Load(context.Background())
	if err != nil {
		t.Fatalf("expected missing state file to load successfully: %v", err)
	}

	if recordID != 0 {
		t.Fatalf("expected record ID 0, got %d", recordID)
	}
}

func TestFileCursorStoreSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cursor.json")

	store, err := NewFileCursorStore(path)
	if err != nil {
		t.Fatalf("expected store creation to succeed: %v", err)
	}

	const expectedRecordID uint64 = 94074458

	if err := store.Save(context.Background(), expectedRecordID); err != nil {
		t.Fatalf("expected save to succeed: %v", err)
	}

	recordID, err := store.Load(context.Background())
	if err != nil {
		t.Fatalf("expected load to succeed: %v", err)
	}

	if recordID != expectedRecordID {
		t.Fatalf(
			"expected record ID %d, got %d",
			expectedRecordID,
			recordID,
		)
	}
}

func TestFileCursorStoreRejectsZeroRecordID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cursor.json")

	store, err := NewFileCursorStore(path)
	if err != nil {
		t.Fatalf("expected store creation to succeed: %v", err)
	}

	err = store.Save(context.Background(), 0)

	if err != ErrInvalidRecordID {
		t.Fatalf(
			"expected ErrInvalidRecordID, got %v",
			err,
		)
	}
}

func TestFileCursorStoreCreatesParentDirectory(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"nested",
		"state",
		"cursor.json",
	)

	store, err := NewFileCursorStore(path)
	if err != nil {
		t.Fatalf("expected store creation to succeed: %v", err)
	}

	if err := store.Save(context.Background(), 123); err != nil {
		t.Fatalf("expected save to create directories: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected state file to exist: %v", err)
	}
}

func TestFileCursorStoreRejectsCorruptedState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cursor.json")

	if err := os.WriteFile(
		path,
		[]byte(`{"lastRecordId":`),
		0600,
	); err != nil {
		t.Fatalf("failed to create corrupted state: %v", err)
	}

	store, err := NewFileCursorStore(path)
	if err != nil {
		t.Fatalf("expected store creation to succeed: %v", err)
	}

	_, err = store.Load(context.Background())

	if err == nil {
		t.Fatal("expected corrupted state to return an error")
	}
}

func TestFileCursorStoreHonorsCanceledContext(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cursor.json")

	store, err := NewFileCursorStore(path)
	if err != nil {
		t.Fatalf("expected store creation to succeed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := store.Load(ctx); err == nil {
		t.Fatal("expected Load to fail with canceled context")
	}

	if err := store.Save(ctx, 123); err == nil {
		t.Fatal("expected Save to fail with canceled context")
	}
}
