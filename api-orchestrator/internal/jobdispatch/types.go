package jobdispatch

type IntegrationVariable struct {
	Name  string `json:"name"`
	Lower string `json:"lower"`
	Upper string `json:"upper"`
}

type WorkerSubmission struct {
	JobID        string                `json:"jobId"`
	FunctionName string                `json:"functionName"`
	Integrand    string                `json:"integrand"`
	Variables    []IntegrationVariable `json:"integrationVariables"`
	Evaluations  int                   `json:"evaluations"`
	Header       string                `json:"header"`
}
