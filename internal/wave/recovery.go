package wave

import (
    "github.com/wyw14/solo-037/internal/model"
)

type ReflectionWindow struct {
    Run model.RunID `json:"run"`
    Revision model.RevisionID `json:"revision"`
    Samples int `json:"samples"`
    Calibrated bool `json:"calibrated"`
}

func NewReflectionWindow(run model.RunID, revision model.RevisionID) ReflectionWindow {
    return ReflectionWindow{Run: run, Revision: revision}
}

func (w ReflectionWindow) Add(revision model.RevisionID, samples int) ReflectionWindow {
    if w.Revision != revision {
        return ReflectionWindow{Run: w.Run, Revision: revision, Samples: samples, Calibrated: samples >= requiredSamples}
    }
    w.Samples += samples
    w.Calibrated = w.Samples >= requiredSamples
    return w
}

const requiredSamples = 3

func (w ReflectionWindow) Valid(revision model.RevisionID) bool {
    return w.Revision == revision && w.Samples >= requiredSamples && w.Calibrated
}
