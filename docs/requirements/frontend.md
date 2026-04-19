# Frontend Requirements (04/2026)

## Purpose

The frontend provides a web-based user interface for defining numerical integration jobs, submitting them to the backend, and inspecting the resulting job status.

The initial version focuses on a simple and functional job submission workflow.

## Initial Scope

The frontend provides a minimal user interface for job submission and for displaying backend responses.

## Core Functional Requirements

- **FE-001** [implemented] [automated-test-covered]  
  The frontend shall allow a user to enter a job name.

- **FE-002** [implemented] [automated-test-covered]  
  The frontend shall allow a user to enter an integrand expression.

- **FE-003** [implemented] [automated-test-covered]  
  The frontend shall allow a user to enter one or more integration variables.

- **FE-004** [implemented] [automated-test-covered]  
  The frontend shall allow a user to specify lower and upper bounds for each integration variable.

- **FE-005** [implemented] [automated-test-covered]  
  The frontend shall allow a user to enter the desired number of evaluations.

- **FE-006** [implemented] [automated-test-covered]  
  The frontend shall allow a user to submit the job to the API orchestrator.

- **FE-007** [implemented] [automated-test-covered]  
  The frontend shall send job data in the JSON structure expected by the backend.

- **FE-008** [implemented] [automated-test-covered]  
  The frontend shall process the backend response and display the returned job ID.

- **FE-009** [implemented] [automated-test-covered]  
  The frontend shall process the backend response and display the returned job status.

- **FE-010** [implemented] [automated-test-covered]  
  The frontend shall display backend-side errors if the submission fails.

## Data Model Requirements

- **FE-011** [implemented] [manually-tested]  
  The frontend shall use request and response types consistent with the backend API.

- **FE-012** [implemented] [automated-test-covered]  
  The frontend shall use a request model for job creation that includes `name`, `integrand`, `variables`, and `evaluations`.

- **FE-013** [implemented] [automated-test-covered]  
  Each variable entry shall include `name`, `lower`, and `upper`.

- **FE-014** [implemented] [automated-test-covered]  
  The frontend shall use a response model for job creation that includes job ID and status.

- **FE-015** [implemented] [automated-test-covered]  
  The frontend shall keep JSON field names aligned with the backend contract.

- **FE-016** [planned] [not-verified]  
  The frontend should later allow a user to upload a JSON file containing a complete job definition.

- **FE-017** [planned] [not-verified]  
  The frontend should later parse such uploaded JSON files and convert them into the same request structure used for manual job submission.

## Initial Non-Functional Requirements

- **FE-018** [implemented] [manually-tested]  
  The frontend shall prioritize simplicity.

- **FE-019** [implemented] [automated-test-covered]  
  The frontend shall prioritize correctness of API interaction.

- **FE-020** [implemented] [automated-test-covered]  
  The frontend shall prioritize clear feedback after submission.

- **FE-021** [implemented] [manually-tested]  
  The frontend shall prioritize ease of local development.

## Current Limitations / Later Extensions

The initial frontend does not yet need to provide:

- a job overview page
- a job detail page
- automatic polling for job state updates
- progress visualization
- convergence visualization
- job cancellation
- authentication or user management
- upload of job definitions from text files in a predefined format

These capabilities may be added in later iterations.