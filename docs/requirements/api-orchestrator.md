# API Orchestrator Requirements (04/2026)

## Purpose

The API orchestrator provides the main backend entry point for creating, validating, storing, and later coordinating numerical integration jobs.

The initial version focuses on accepting job submissions, validating input, persisting accepted jobs, and exposing basic HTTP endpoints.

## Initial Scope

The API orchestrator provides:

- an HTTP API for frontend-backend communication
- a health endpoint
- a job creation endpoint
- request validation for submitted jobs
- persistence of accepted jobs
- a foundation for later worker coordination and job status retrieval

## Core Functional Requirements

- **API-001** [implemented] [manually-tested]  
  The API orchestrator shall accept HTTP requests from the frontend.

- **API-002** [implemented] [manually-tested]  
  The API orchestrator shall accept job creation requests in JSON format.

- **API-003** [implemented] [manually-tested]  
  The API orchestrator shall validate the request structure before accepting a job.

- **API-004** [implemented] [manually-tested]  
  The API orchestrator shall reject malformed JSON requests.

- **API-005** [implemented] [manually-tested]  
  The API orchestrator shall reject unknown JSON fields.

- **API-006** [implemented] [manually-tested]  
  The API orchestrator shall reject request bodies that contain trailing data or multiple JSON objects.

- **API-007** [implemented] [automated-test-covered]  
  The API orchestrator shall validate required job fields before acceptance.

- **API-008** [implemented] [manually-tested]  
  The API orchestrator shall assign or create a unique job identifier for each accepted job.

- **API-009** [implemented] [manually-tested]  
  The API orchestrator shall persist accepted jobs in the database.

- **API-010** [implemented] [manually-tested]  
  The API orchestrator shall return a structured JSON response for accepted jobs.

- **API-011** [implemented] [manually-tested]  
  The API orchestrator shall return structured JSON error responses for rejected requests.

- **API-012** [implemented] [manually-tested]  
  The API orchestrator shall allow a client to submit a job name.

- **API-013** [implemented] [manually-tested]  
  The API orchestrator shall allow a client to submit an integrand expression.

- **API-014** [implemented] [manually-tested]  
  The API orchestrator shall allow a client to submit one or more integration variables.

- **API-015** [implemented] [manually-tested]  
  The API orchestrator shall allow a client to submit lower and upper bounds per variable.

- **API-016** [implemented] [manually-tested]  
  The API orchestrator shall allow a client to submit a number of evaluations.

- **API-017** [implemented] [automated-test-covered]  
  The API orchestrator shall validate that the job name is not empty.

- **API-018** [implemented] [automated-test-covered]  
  The API orchestrator shall validate that the integrand is not empty.

- **API-019** [implemented] [automated-test-covered]  
  The API orchestrator shall validate that at least one integration variable is provided.

- **API-020** [implemented] [automated-test-covered]  
  The API orchestrator shall validate that the evaluation count is greater than zero.

- **API-021** [implemented] [automated-test-covered]  
  The API orchestrator shall validate that each variable definition is structurally valid.

- **API-022** [implemented] [automated-test-covered]  
  The API orchestrator shall validate integrand semantics for valid expressions.

- **API-023** [implemented] [automated-test-covered]  
  The API orchestrator shall reject unknown functions during semantic validation.

- **API-024** [implemented] [automated-test-covered]  
  The API orchestrator shall reject unknown identifiers during semantic validation.

- **API-025** [implemented] [automated-test-covered]  
  The API orchestrator shall support built-in constants during semantic validation.

- **API-026** [implemented] [automated-test-covered]  
  The API orchestrator shall support lexical and syntactic parsing of integrand expressions.

- **API-027** [implemented] [automated-test-covered]  
  The API orchestrator shall reject malformed integrand expressions during parsing.

- **API-028** [planned] [not-verified]  
  The API orchestrator should later provide a job status endpoint.

- **API-029** [planned] [not-verified]  
  The API orchestrator should later provide a job list endpoint.

- **API-030** [planned] [not-verified]  
  The API orchestrator should later provide a job details endpoint.

- **API-031** [planned] [not-verified]  
  The API orchestrator should later provide a job cancellation endpoint.

- **API-032** [planned] [not-verified]  
  The API orchestrator should later coordinate with the compute worker.

- **API-033** [planned] [not-verified]  
  The API orchestrator should later provide progress and result retrieval endpoints.

## API Requirements

- **API-034** [implemented] [manually-tested]  
  The API orchestrator shall expose a health endpoint for basic service checks.

- **API-035** [implemented] [manually-tested]  
  The API orchestrator shall expose a job creation endpoint.

- **API-036** [implemented] [manually-tested]  
  The API orchestrator shall use JSON as the request and response format.

- **API-037** [implemented] [manually-tested]  
  The API orchestrator shall return response fields that are stable and consistent with the frontend data model.

- **API-038** [implemented] [manually-tested]  
  The API orchestrator shall return an accepted job response containing at least the job ID and the job status.

- **API-039** [implemented] [manually-tested]  
  The API orchestrator shall return validation or decoding failures as structured error responses.

## Persistence Requirements

- **API-040** [implemented] [manually-tested]  
  The API orchestrator shall persist accepted jobs through a persistence layer.

- **API-041** [implemented] [manually-tested]  
  The API orchestrator shall not require the frontend to communicate directly with the database.

- **API-042** [implemented] [manually-tested]  
  The API orchestrator shall store enough information to reconstruct a submitted job and its current state.

## Initial Non-Functional Requirements

- **API-046** [implemented] [manually-tested]  
  The API orchestrator shall prioritize simplicity.

- **API-047** [implemented] [automated-test-covered]  
  The API orchestrator shall prioritize clear request validation.

- **API-048** [implemented] [manually-tested]  
  The API orchestrator shall prioritize stable API behavior.

- **API-049** [implemented] [manually-tested]  
  The API orchestrator shall prioritize debuggability.

- **API-050** [implemented] [manually-tested]  
  The API orchestrator shall prioritize ease of local development.

## Current Limitations / Later Extensions

The initial API orchestrator does not yet need to provide:

- authentication or authorization
- multi-user support
- advanced job scheduling
- direct distributed worker orchestration
- progress streaming
- result streaming
- production-grade scalability features

These capabilities may be added in later iterations.