package jobdispatch

import "github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"

type WorkerSubmission struct {
	JobID        string             `json:"jobId"`
	FunctionName string             `json:"functionName"`
	Integrand    string             `json:"integrand"`
	Variables    []job.VariableSpec `json:"integrationVariables"`
	Evaluations  int                `json:"evaluations"`
	Header       string             `json:"header"`
}
