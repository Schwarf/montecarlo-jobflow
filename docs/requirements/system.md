# System-Level Requirements (04/2026)

## Purpose

The system shall allow a user to define numerical integration jobs, submit them, execute them asynchronously, persist their state, and inspect job progress and results through a web UI.

The first version shall prefer simplicity over scale.

## Initial Scope

The system shall provide:

- a web frontend for job submission and status inspection
- an API orchestrator as the main entry point
- a persistence layer for storing jobs and their state
- a compute worker for executing numerical integration jobs

## Core Functional Requirements

The system shall allow a user to:

- define a numerical integration job
- specify an integrand and integration variables with bounds
- submit a job through the frontend
- have the job validated before acceptance
- have accepted jobs persisted
- have jobs executed asynchronously
- inspect job state and later the final result

The system should later allow a user to:

- inspect progress during execution
- inspect convergence-related information
- cancel running or queued jobs

## High-Level Architectural Requirements

The system shall distinguish the concerns of:

- user interaction
- API orchestration
- persistence
- numerical computation
- later analysis

In the initial version, some of these concerns may still be implemented within the same deployable component, 
as long as their responsibilities are kept logically separated in the code.

The system shall be modular such that components can evolve independently.

The system shall run locally in an initial development setup.

The system should later be extensible to a more distributed architecture.

## Initial Non-Functional Requirements

The first version shall prioritize:

- simplicity
- clear component boundaries
- debuggability
- ease of local development