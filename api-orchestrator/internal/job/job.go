package job

import "time"

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusCanceled  Status = "canceled"
)

type VariableSpec struct {
	Name  string `json:"name"`
	Lower string `json:"lower"`
	Upper string `json:"upper"`
}

type Job struct {
	ID                   string
	Name                 string
	Integrand            string
	IntegrationVariables []VariableSpec
	Evaluations          int
	Status               Status
	ErrorMessage         *string
	ResultJSON           *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
