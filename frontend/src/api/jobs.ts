import type { CreateJobRequest, CreateJobResponse } from "../types/job";

export async function createJob(
    request: CreateJobRequest
): Promise<CreateJobResponse> {
    const response = await fetch("http://localhost:8080/api/v1/jobs", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
    });

    if (!response.ok) {
        throw new Error(`request failed with status ${response.status}`);
    }

    return response.json() as Promise<CreateJobResponse>;
}