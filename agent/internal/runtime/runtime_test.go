package runtime

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/Anis-hacker26/SignalLens/agent/internal/evidence"
	"github.com/Anis-hacker26/SignalLens/agent/internal/pipeline"
	"github.com/Anis-hacker26/SignalLens/agent/internal/state"
	"github.com/Anis-hacker26/SignalLens/agent/internal/telemetry"
)

type fakeCollector struct {
	events      []telemetry.RawEvent
	err         error
	cursor      uint64
	commitCount int
}

func (f *fakeCollector) Collect(
	ctx context.Context,
) ([]telemetry.RawEvent, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.events, nil
}

func (f *fakeCollector) SetCursor(recordID uint64) {
	f.cursor = recordID
}

func (f *fakeCollector) Commit(recordID uint64) {
	if recordID > f.cursor {
		f.cursor = recordID
	}

	f.commitCount++
}

type fakeCursorStore struct {
	recordID uint64
	err      error
}

func (f *fakeCursorStore) Load(ctx context.Context) (uint64, error) {
	if f.err != nil {
		return 0, f.err
	}

	return f.recordID, nil
}

func (f *fakeCursorStore) Save(ctx context.Context, recordID uint64) error {
	if f.err != nil {
		return f.err
	}

	f.recordID = recordID
	return nil
}

var _ state.CursorStore = (*fakeCursorStore)(nil)

func TestRuntimeRunProcessesCollectedEvents(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{
		events: []telemetry.RawEvent{
			{
				EventID:   1,
				Timestamp: time.Date(2026, 11, 13, 1, 0, 0, 0, time.UTC),
				Provider:  "Microsoft-Windows-Sysmon",
				Computer:  "TEST-HOST",
				Message:   "<Event>process</Event>",
			},
			{
				EventID:   3,
				Timestamp: time.Date(2026, 11, 13, 1, 1, 0, 0, time.UTC),
				Provider:  "Microsoft-Windows-Sysmon",
				Computer:  "TEST-HOST",
				Message:   "<Event>network</Event>",
			},
		},
	}

	var callbackResults []evidence.Evidence

	runtime, err := New(
		eventCollector,
		eventPipeline,
		&fakeCursorStore{},
		time.Second,
		func(_ context.Context, results []evidence.Evidence) error {
			callbackResults = results
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	results, err := runtime.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf(
			"expected 2 evidence records, got %d",
			len(results),
		)
	}

	if results[0].Event.Type != "process" {
		t.Fatalf(
			"expected first event to be process, got %q",
			results[0].Event.Type,
		)
	}

	if results[1].Event.Type != "network" {
		t.Fatalf(
			"expected second event to be network, got %q",
			results[1].Event.Type,
		)
	}

	if len(callbackResults) != 2 {
		t.Fatalf(
			"expected callback to receive 2 evidence records, got %d",
			len(callbackResults),
		)
	}

	if callbackResults[0].ID != results[0].ID {
		t.Fatalf("expected callback evidence to match runtime result")
	}

	if callbackResults[1].ID != results[1].ID {
		t.Fatalf("expected callback evidence to match runtime result")
	}
}

func TestRuntimeRunReturnsEmptyEvidenceWhenNoEvents(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{}

	var callbackCalled bool
	var callbackResults []evidence.Evidence

	runtime, err := New(
		eventCollector,
		eventPipeline,
		&fakeCursorStore{},
		time.Second,
		func(_ context.Context, results []evidence.Evidence) error {
			callbackCalled = true
			callbackResults = results
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	results, err := runtime.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if results == nil {
		t.Fatal("expected non-nil empty evidence slice")
	}

	if len(results) != 0 {
		t.Fatalf(
			"expected 0 evidence records, got %d",
			len(results),
		)
	}

	if !callbackCalled {
		t.Fatal("expected callback to be called")
	}

	if callbackResults == nil {
		t.Fatal("expected callback to receive non-nil empty evidence slice")
	}

	if len(callbackResults) != 0 {
		t.Fatalf(
			"expected callback to receive 0 evidence records, got %d",
			len(callbackResults),
		)
	}
}

func TestRuntimeRejectsNilCollector(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	_, err = New(
		nil,
		eventPipeline,
		&fakeCursorStore{},
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err == nil {
		t.Fatal("expected error for nil collector")
	}
}

func TestRuntimeRejectsNilPipeline(t *testing.T) {
	eventCollector := &fakeCollector{}

	_, err := New(
		eventCollector,
		nil,
		&fakeCursorStore{},
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err == nil {
		t.Fatal("expected error for nil pipeline")
	}
}

func TestRuntimeRejectsInvalidInterval(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{}

	_, err = New(
		eventCollector,
		eventPipeline,
		&fakeCursorStore{},
		0,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestRuntimeRejectsNilContext(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{}

	runtime, err := New(
		eventCollector,
		eventPipeline,
		&fakeCursorStore{},
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	_, err = runtime.Run(nil)
	if err == nil {
		t.Fatal("expected error for nil context")
	}
}

func TestRuntimeStartAndStop(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{
		events: []telemetry.RawEvent{
			{
				EventID:   1,
				Timestamp: time.Now(),
				Provider:  "Microsoft-Windows-Sysmon",
				Computer:  "TEST-HOST",
				Message:   "<Event>process</Event>",
			},
		},
	}

	runtime, err := New(
		eventCollector,
		eventPipeline,
		&fakeCursorStore{},
		10*time.Millisecond,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	ctx := context.Background()

	if err := runtime.Start(ctx); err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}

	if err := runtime.Start(ctx); err == nil {
		t.Fatal("expected error when starting runtime twice")
	}

	stopCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	if err := runtime.Stop(stopCtx); err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}

	if err := runtime.Stop(stopCtx); err != nil {
		t.Fatalf("expected repeated stop to be safe, got: %v", err)
	}
}

func TestRuntimeStartRejectsNilContext(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{}

	runtime, err := New(
		eventCollector,
		eventPipeline,
		&fakeCursorStore{},
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	if err := runtime.Start(nil); err == nil {
		t.Fatal("expected error for nil context")
	}
}

func TestRuntimeStopRejectsNilContext(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{}

	runtime, err := New(
		eventCollector,
		eventPipeline,
		&fakeCursorStore{},
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	if err := runtime.Stop(nil); err == nil {
		t.Fatal("expected error for nil context")
	}
}

func TestRuntimeRunReturnsCollectorError(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	expectedErr := errors.New("collector failed")

	eventCollector := &fakeCollector{
		err: expectedErr,
	}

	runtime, err := New(
		eventCollector,
		eventPipeline,
		&fakeCursorStore{},
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	_, err = runtime.Run(context.Background())
	if err == nil {
		t.Fatal("expected collector error")
	}
}

func TestRuntimeRunDoesNotCheckpointWhenDeliveryFails(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{
		events: []telemetry.RawEvent{
			{
				RecordID: 42,
				EventID:  1,
				Timestamp: time.Date(
					2026,
					11,
					13,
					1,
					0,
					0,
					0,
					time.UTC,
				),
				Provider: "Microsoft-Windows-Sysmon",
				Computer: "TEST-HOST",
				Message:  "<Event>process</Event>",
			},
		},
	}

	cursorStore := &fakeCursorStore{
		recordID: 10,
	}

	expectedErr := errors.New("delivery failed")

	runtime, err := New(
		eventCollector,
		eventPipeline,
		cursorStore,
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return expectedErr
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	_, err = runtime.Run(context.Background())
	if err == nil {
		t.Fatal("expected delivery error")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected delivery error %v, got %v",
			expectedErr,
			err,
		)
	}

	if cursorStore.recordID != 10 {
		t.Fatalf(
			"expected cursor to remain 10, got %d",
			cursorStore.recordID,
		)
	}

	if eventCollector.cursor != 0 {
		t.Fatalf(
			"expected collector cursor to remain 0, got %d",
			eventCollector.cursor,
		)
	}

	if eventCollector.commitCount != 0 {
		t.Fatalf(
			"expected 0 cursor commits, got %d",
			eventCollector.commitCount,
		)
	}
}

func TestRuntimeRunCheckpointsAfterSuccessfulDelivery(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{
		events: []telemetry.RawEvent{
			{
				RecordID: 42,
				EventID:  1,
				Timestamp: time.Date(
					2026,
					11,
					13,
					1,
					0,
					0,
					0,
					time.UTC,
				),
				Provider: "Microsoft-Windows-Sysmon",
				Computer: "TEST-HOST",
				Message:  "<Event>process</Event>",
			},
		},
	}

	cursorStore := &fakeCursorStore{
		recordID: 10,
	}

	runtime, err := New(
		eventCollector,
		eventPipeline,
		cursorStore,
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	_, err = runtime.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cursorStore.recordID != 42 {
		t.Fatalf(
			"expected persisted cursor to be 42, got %d",
			cursorStore.recordID,
		)
	}

	if eventCollector.cursor != 42 {
		t.Fatalf(
			"expected collector cursor to be 42, got %d",
			eventCollector.cursor,
		)
	}

	if eventCollector.commitCount != 1 {
		t.Fatalf(
			"expected 1 cursor commit, got %d",
			eventCollector.commitCount,
		)
	}
}
func TestRuntimeStartLoadsPersistedCursor(t *testing.T) {
	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("unexpected error creating pipeline: %v", err)
	}

	eventCollector := &fakeCollector{}

	cursorStore := &fakeCursorStore{
		recordID: 42,
	}

	runtime, err := New(
		eventCollector,
		eventPipeline,
		cursorStore,
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating runtime: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := runtime.Start(ctx); err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}

	if eventCollector.cursor != 42 {
		t.Fatalf(
			"expected collector cursor to be initialized to 42, got %d",
			eventCollector.cursor,
		)
	}

	stopCtx, stopCancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer stopCancel()

	if err := runtime.Stop(stopCtx); err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}
}
func TestRuntimePersistsCursorAcrossRestart(t *testing.T) {
	tempDir := t.TempDir()
	statePath := filepath.Join(
		tempDir,
		"signallens-state.json",
	)

	// Simulate the first agent process.
	cursorStore1, err := state.NewFileCursorStore(statePath)
	if err != nil {
		t.Fatalf("failed to create first cursor store: %v", err)
	}

	eventCollector1 := &fakeCollector{
		events: []telemetry.RawEvent{
			{
				RecordID: 42,
				Provider: "Microsoft-Windows-Sysmon",
				EventID:  1,
				Timestamp: time.Date(
					2026,
					11,
					13,
					1,
					0,
					0,
					0,
					time.UTC,
				),
				Computer: "TEST-HOST",
				Message:  "<Event>process</Event>",
			},
		},
	}

	eventPipeline, err := pipeline.NewDefault()
	if err != nil {
		t.Fatalf("failed to create pipeline: %v", err)
	}

	runtime1, err := New(
		eventCollector1,
		eventPipeline,
		cursorStore1,
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("failed to create first runtime: %v", err)
	}

	if _, err := runtime1.Run(context.Background()); err != nil {
		t.Fatalf("first runtime run failed: %v", err)
	}

	if eventCollector1.cursor != 42 {
		t.Fatalf(
			"expected first collector cursor 42, got %d",
			eventCollector1.cursor,
		)
	}

	// Simulate the agent process restarting.
	cursorStore2, err := state.NewFileCursorStore(statePath)
	if err != nil {
		t.Fatalf("failed to create second cursor store: %v", err)
	}

	eventCollector2 := &fakeCollector{}

	runtime2, err := New(
		eventCollector2,
		eventPipeline,
		cursorStore2,
		time.Second,
		func(context.Context, []evidence.Evidence) error {
			return nil
		},
	)
	if err != nil {
		t.Fatalf("failed to create second runtime: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := runtime2.Start(ctx); err != nil {
		t.Fatalf("failed to start second runtime: %v", err)
	}

	if eventCollector2.cursor != 42 {
		t.Fatalf(
			"expected restarted collector cursor 42, got %d",
			eventCollector2.cursor,
		)
	}

	stopCtx, stopCancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer stopCancel()

	if err := runtime2.Stop(stopCtx); err != nil {
		t.Fatalf("failed to stop second runtime: %v", err)
	}
}
