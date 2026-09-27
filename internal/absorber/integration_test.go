package absorber

import (
    "testing"
    "time"

    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/tank"
)

func evidenceBatch(run model.RunID, revision model.RevisionID, source string, count int, at time.Time) []model.Evidence {
    samples := make([]model.Evidence, 0, count)
    for index := 0; index < count; index++ {
        samples = append(samples, model.NewEvidence(run, revision, source, model.EvidenceWaveTime, 0.1, at.Add(time.Duration(index)*time.Millisecond)))
    }
    return samples
}

func TestCalibrateFromProbeIgnoresOldRevisionEvidence(t *testing.T) {
    run := model.RunID("run-1")
    oldRevision := model.RevisionID("rev-old")
    newRevision := model.RevisionID("rev-new")
    at := time.Unix(1700000000, 0)
    controller := NewController(tank.NewBoundary())

    oldBatch := evidenceBatch(run, oldRevision, "probe-a", 3, at)
    if !CalibrateFromProbe(controller, oldBatch, oldRevision) {
        t.Fatal("old revision should calibrate with its own sufficient batch")
    }

    // Two new-revision samples mixed with three stale old-revision samples:
    // the new revision is still short and must not borrow the old evidence.
    mixed := append(oldBatch, evidenceBatch(run, newRevision, "probe-a", 2, at.Add(10*time.Millisecond))...)
    if CalibrateFromProbe(controller, mixed, newRevision) {
        t.Fatal("new revision must not calibrate using old revision evidence")
    }
    if controller.ReflectionReady(run, newRevision) {
        t.Fatal("new revision must not be reported ready below the threshold")
    }
    if controller.ReflectionReady(run, oldRevision) {
        t.Fatal("old revision must no longer be valid after the window switched")
    }

    // A single sufficient batch for the new revision takes effect immediately,
    // even on a fresh window that has never seen any prior batch.
    fresh := NewController(tank.NewBoundary())
    newBatch := evidenceBatch(run, newRevision, "probe-a", 3, at.Add(20*time.Millisecond))
    if !CalibrateFromProbe(fresh, newBatch, newRevision) {
        t.Fatal("new revision should calibrate on the first sufficient batch")
    }
    if !fresh.ReflectionReady(run, newRevision) {
        t.Fatal("new revision must be reported ready in the same submission")
    }
}
