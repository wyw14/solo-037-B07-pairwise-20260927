package wave

import (
    "testing"

    "github.com/wyw14/solo-037/internal/model"
)

func TestReflectionWindowCalibratesWithinFirstBatchAfterRevisionSwitch(t *testing.T) {
    run := model.RunID("run-1")
    oldRevision := model.RevisionID("rev-old")
    newRevision := model.RevisionID("rev-new")

    window := NewReflectionWindow(run, oldRevision)
    window = window.Add(oldRevision, 3)
    if !window.Valid(oldRevision) {
        t.Fatal("old revision should be calibrated after sufficient samples")
    }

    window = window.Add(newRevision, 3)
    if window.Samples != 3 {
        t.Fatalf("new revision must not inherit old sample count: got %d", window.Samples)
    }
    if !window.Valid(newRevision) {
        t.Fatal("new revision with a sufficient single batch must be calibrated immediately")
    }
    if window.Valid(oldRevision) {
        t.Fatal("old revision evidence must not remain valid after switch")
    }
}

func TestReflectionWindowStaysUncalibratedBelowThreshold(t *testing.T) {
    run := model.RunID("run-2")
    revision := model.RevisionID("rev-1")
    window := NewReflectionWindow(run, revision).Add(revision, 2)
    if window.Valid(revision) {
        t.Fatal("window must not be calibrated below the sample threshold")
    }
    window = window.Add(revision, 1)
    if !window.Valid(revision) {
        t.Fatal("window must calibrate once the threshold is reached")
    }
}
