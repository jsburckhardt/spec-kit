// Package config provides configuration structures and constants for gospecify
package config

import (
	"sync"
	"time"
)

// ProjectConfig holds the configuration for a project initialization
type ProjectConfig struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	AIAssistant string    `json:"ai_assistant"`
	ScriptType  string    `json:"script_type"`
	NoGit       bool      `json:"no_git"`
	Force       bool      `json:"force"`
	IgnoreTools bool      `json:"ignore_tools"`
	SkipTLS     bool      `json:"skip_tls"`
	Debug       bool      `json:"debug"`
	GitHubToken string    `json:"github_token,omitempty"`
	Here        bool      `json:"here"`
	CreatedAt   time.Time `json:"created_at"`
}

// StepTracker manages hierarchical progress tracking with live updates
type StepTracker struct {
	Title       string         `json:"title"`
	Steps       []Step         `json:"steps"`
	StatusOrder map[string]int `json:"-"`
	refreshCb   func()         `json:"-"`
	mu          sync.RWMutex   `json:"-"`
}

// Step represents a single step in the progress tracking
type Step struct {
	Key     string    `json:"key"`
	Label   string    `json:"label"`
	Status  Status    `json:"status"`
	Detail  string    `json:"detail"`
	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
}

// Status represents the status of a step
type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusError   Status = "error"
	StatusSkipped Status = "skipped"
)

// AttachRefresh attaches a callback function for live UI updates
func (st *StepTracker) AttachRefresh(cb func()) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.refreshCb = cb
}

// Add adds a new step to the tracker
func (st *StepTracker) Add(key, label string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.StatusOrder == nil {
		st.StatusOrder = map[string]int{
			"pending": 0,
			"running": 1,
			"done":    2,
			"error":   3,
			"skipped": 4,
		}
	}

	// Check if step already exists
	for _, step := range st.Steps {
		if step.Key == key {
			return
		}
	}

	st.Steps = append(st.Steps, Step{
		Key:    key,
		Label:  label,
		Status: StatusPending,
	})

	st.maybeRefresh()
}

// Start marks a step as running
func (st *StepTracker) Start(key, detail string) {
	st.update(key, StatusRunning, detail)
}

// Complete marks a step as done
func (st *StepTracker) Complete(key, detail string) {
	st.update(key, StatusDone, detail)
}

// Error marks a step as error
func (st *StepTracker) Error(key, detail string) {
	st.update(key, StatusError, detail)
}

// Skip marks a step as skipped
func (st *StepTracker) Skip(key, detail string) {
	st.update(key, StatusSkipped, detail)
}

// update updates a step's status and detail
func (st *StepTracker) update(key string, status Status, detail string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	for i := range st.Steps {
		if st.Steps[i].Key == key {
			st.Steps[i].Status = status
			if detail != "" {
				st.Steps[i].Detail = detail
			}
			if status == StatusRunning && st.Steps[i].Started.IsZero() {
				st.Steps[i].Started = time.Now()
			}
			if (status == StatusDone || status == StatusError || status == StatusSkipped) && st.Steps[i].Ended.IsZero() {
				st.Steps[i].Ended = time.Now()
			}
			st.maybeRefresh()
			return
		}
	}

	// If not found, add it
	st.Steps = append(st.Steps, Step{
		Key:    key,
		Label:  key,
		Status: status,
		Detail: detail,
	})
	st.maybeRefresh()
}

// maybeRefresh calls the refresh callback if set
func (st *StepTracker) maybeRefresh() {
	if st.refreshCb != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Ignore panics in refresh callback
					_ = r
				}
			}()
			st.refreshCb()
		}()
	}
}

// GetSteps returns a copy of the current steps for safe reading
func (st *StepTracker) GetSteps() []Step {
	st.mu.RLock()
	defer st.mu.RUnlock()

	steps := make([]Step, len(st.Steps))
	copy(steps, st.Steps)
	return steps
}

// GetTitle returns the tracker title
func (st *StepTracker) GetTitle() string {
	return st.Title
}
