# Frontend Requirements (04/2026)

## Purpose

The frontend shall provide a web-based user interface for defining numerical integration jobs, submitting them to the backend, and inspecting the resulting job status.

The initial version shall focus on a simple and functional job submission workflow.

## Initial Scope

The frontend shall provide:

- a form for defining and submitting integration jobs
- a way to enter the job name
- a way to enter the integrand
- a way to enter integration variables with lower and upper bounds
- a way to enter the number of evaluations
- a way to display the backend response after job submission
- a way to display validation or submission errors

## Core Functional Requirements

The frontend shall allow a user to:

- enter a job name
- enter an integrand expression
- enter one or more integration variables
- specify lower and upper bounds for each integration variable
- enter the desired number of evaluations
- submit the job to the API orchestrator

The frontend shall send job data in the JSON structure expected by the backend.

The frontend shall process the backend response and display at least:

- the returned job ID
- the returned job status

The frontend shall display backend-side errors if the submission fails.

## Data Model Requirements

The frontend shall use request and response types consistent with the backend API.

The frontend shall use a request model for job creation that includes:

- name
- integrand
- variables
- evaluations

Each variable entry shall include:

- name
- lower
- upper

The frontend shall use a response model for job creation that includes:

- job ID
- status

The frontend shall keep JSON field names aligned with the backend contract.

The frontend shall display backend-side errors if the submission fails.

The frontend should later allow a user to upload a text file containing a complete job definition in a predefined format.

The frontend should later parse such uploaded files and convert them into the same request structure used for manual job submission.

And under later extensions:
## Initial Non-Functional Requirements

The frontend shall prioritize:

- simplicity
- correctness of API interaction
- clear feedback after submission
- ease of local development

The frontend shall be implemented in TypeScript.

The frontend should keep API interaction logic and UI logic reasonably separated.

The frontend should remain easy to extend with later views such as:

- job list view
- job detail view
- progress display
- result display

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