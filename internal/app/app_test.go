package app

import (
	"testing"

	"cfvgtuai/internal/codeforces"
)

func TestTracker(t *testing.T) {
	tracker := NewTracker()

	submission1 := codeforces.Submission{
		ID:      1,
		Problem: codeforces.Problem{Name: "ProblemA"},
		Verdict: "WRONG_ANSWER",
	}
	submission2 := codeforces.Submission{
		ID:      2,
		Problem: codeforces.Problem{Name: "ProblemA"},
		Verdict: "OK",
	}

	// Process first submission (wrong answer)
	problemID1 := submission1.Problem.Name
	if submission1.Verdict != "OK" {
		tracker.wrongAttempts[problemID1]++
	}
	tracker.processedSubmissions[submission1.ID] = true

	if tracker.wrongAttempts[problemID1] != 1 {
		t.Errorf("Expected wrong attempts to be 1, but got %d", tracker.wrongAttempts[problemID1])
	}

	// Process second submission (correct answer)
	problemID2 := submission2.Problem.Name
	if submission2.Verdict == "OK" {
		if !tracker.processedSubmissions[submission2.ID] {
			wrongAttempts := tracker.wrongAttempts[problemID2]
			if wrongAttempts != 1 {
				t.Errorf("Expected wrong attempts to be 1, but got %d", wrongAttempts)
			}
			tracker.processedSubmissions[submission2.ID] = true
			delete(tracker.wrongAttempts, problemID2)
		}
	}

	if _, exists := tracker.wrongAttempts[problemID2]; exists {
		t.Errorf("Expected wrong attempts for %s to be deleted, but it still exists", problemID2)
	}

	if !tracker.processedSubmissions[submission1.ID] || !tracker.processedSubmissions[submission2.ID] {
		t.Errorf("Expected both submissions to be marked as processed")
	}
}
