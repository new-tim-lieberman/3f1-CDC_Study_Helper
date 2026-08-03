package quiz

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// QuestionStats tracks one question's progress through the Leitner boxes.
type QuestionStats struct {
	Box    int       `json:"box"`
	Due    time.Time `json:"due"`
	Seen   int       `json:"seen"`
	Missed int       `json:"missed"`
}

// Stats maps a question key (see QuestionKey) to its Leitner progress.
type Stats map[string]QuestionStats

// boxIntervals[i] is how long a question sits before it's due again once it
// reaches box i. Box 0 is always due immediately (new or just-missed).
var boxIntervals = []time.Duration{
	0,
	24 * time.Hour,
	3 * 24 * time.Hour,
	7 * 24 * time.Hour,
	16 * 24 * time.Hour,
}

var maxBoxIndex = len(boxIntervals) - 1

// StatsPath returns the default location for the persisted stats file.
func StatsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".3f1_cdc_stats.json"), nil
}

// LoadStats reads the stats file at path, returning an empty Stats if it
// doesn't exist yet.
func LoadStats(path string) (Stats, error) {
	stats := make(Stats)
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return stats, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &stats); err != nil {
		return nil, fmt.Errorf("parsing stats file: %w", err)
	}
	return stats, nil
}

// SaveStats writes stats to path.
func SaveStats(path string, stats Stats) error {
	raw, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

// QuestionKey derives a stable identity for q from its module and text, so
// progress survives across runs even though questions carry no explicit ID.
func QuestionKey(q Question) string {
	h := sha1.Sum([]byte(fmt.Sprintf("%d:%s", q.ModuleID, q.Question)))
	return hex.EncodeToString(h[:])
}

// Record updates q's Leitner box after an attempt: a correct answer promotes
// it to the next box (reviewed further out), a miss drops it back to box 0.
func (s Stats) Record(q Question, correct bool) {
	key := QuestionKey(q)
	st := s[key]
	st.Seen++
	if correct {
		if st.Box < maxBoxIndex {
			st.Box++
		}
	} else {
		st.Box = 0
		st.Missed++
	}
	st.Due = time.Now().Add(boxIntervals[st.Box])
	s[key] = st
}

// Due reports whether q has never been studied or its review interval has
// elapsed.
func (s Stats) Due(q Question) bool {
	st, ok := s[QuestionKey(q)]
	if !ok {
		return true
	}
	return !time.Now().Before(st.Due)
}
