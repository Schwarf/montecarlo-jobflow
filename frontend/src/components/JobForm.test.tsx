import {describe, expect, it, vi, beforeEach} from "vitest";
import {fireEvent, render, screen} from "@testing-library/react";
import {JobForm} from "./JobForm";
import * as jobsApi from "../api/jobs";


vi.mock("../api/jobs", () => ({
    createJob: vi.fn(),
}));

describe("JobForm", () => {

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders the main form fields", () => {
        render(<JobForm/>);

        expect(screen.getByLabelText(/job name/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/evaluations/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/integrand/i)).toBeInTheDocument();
        expect(screen.getByRole("button", {name: /create job/i})).toBeInTheDocument();
    });

    it("adds another integration variable", () => {
        render(<JobForm/>);

        expect(screen.getAllByLabelText(/^name$/i)).toHaveLength(1);

        fireEvent.click(screen.getByRole("button", {name: /add variable/i}));

        expect(screen.getAllByLabelText(/^name$/i)).toHaveLength(2);
    });

    it("removes an integration variable", () => {
        render(<JobForm/>);

        fireEvent.click(screen.getByRole("button", {name: /add variable/i}));
        expect(screen.getAllByLabelText(/^name$/i)).toHaveLength(2);

        const removeButtons = screen.getAllByRole("button", {name: /remove this variable/i});
        fireEvent.click(removeButtons[1]);

        expect(screen.getAllByLabelText(/^name$/i)).toHaveLength(1);
    });

    it("submits the form and shows the success message", async () => {
        vi.mocked(jobsApi.createJob).mockResolvedValue({
            jobId: "abc-123",
            status: "queued",
        });

        render(<JobForm/>);

        fireEvent.change(screen.getByLabelText(/job name/i), {
            target: {value: "test-job"},
        });

        fireEvent.change(screen.getByLabelText(/integrand/i), {
            target: {value: "x^2"},
        });

        fireEvent.click(screen.getByRole("button", {name: /create job/i}));

        expect(await screen.findByText(/job created: abc-123 \(queued\)/i)).toBeInTheDocument();
        expect(jobsApi.createJob).toHaveBeenCalled();
    });

    it("shows an error message when submission fails", async () => {
        vi.mocked(jobsApi.createJob).mockRejectedValue(new Error("request failed"));

        render(<JobForm/>);

        fireEvent.change(screen.getByLabelText(/job name/i), {
            target: {value: "test-job"},
        });

        fireEvent.change(screen.getByLabelText(/integrand/i), {
            target: {value: "x^2"},
        });

        fireEvent.click(screen.getByRole("button", {name: /create job/i}));

        expect(await screen.findByText(/Failed to create job/i)).toBeInTheDocument();
    });

    it("submits the expected job payload", async () => {
        vi.mocked(jobsApi.createJob).mockResolvedValue({
            jobId: "abc-123",
            status: "queued",
        });

        render(<JobForm/>);

        fireEvent.change(screen.getByLabelText(/job name/i), {
            target: {value: "test-job"},
        });

        fireEvent.change(screen.getByLabelText(/integrand/i), {
            target: {value: "(1+x^2)^4"},
        });

        fireEvent.change(screen.getByLabelText(/evaluations/i), {
            target: {value: "5000"},
        });

        fireEvent.change(screen.getByLabelText(/^name$/i), {
            target: {value: "x"},
        });

        fireEvent.change(screen.getByLabelText(/lower/i), {
            target: {value: "0"},
        });

        fireEvent.change(screen.getByLabelText(/upper/i), {
            target: {value: "1"},
        });

        fireEvent.click(screen.getByRole("button", {name: /create job/i}));

        await screen.findByText(/job created: abc-123 \(queued\)/i);

        expect(jobsApi.createJob).toHaveBeenCalledWith({
            name: "test-job",
            integrand: "(1+x^2)^4",
            variables: [
                {
                    name: "x",
                    lower: "0",
                    upper: "1",
                },
            ],
            evaluations: 5000,
        });
    });

    it("disables removing the last remaining integration variable", () => {
        render(<JobForm />);

        const removeButton = screen.getByRole("button", { name: /remove this variable/i });
        expect(removeButton).toBeDisabled();
    });

    it("submits multiple integration variables in the payload", async () => {
        vi.mocked(jobsApi.createJob).mockResolvedValue({
            jobId: "abc-123",
            status: "queued",
        });

        render(<JobForm />);

        fireEvent.change(screen.getByLabelText(/job name/i), {
            target: { value: "test-job" },
        });

        fireEvent.change(screen.getByLabelText(/integrand/i), {
            target: { value: "x+y" },
        });

        fireEvent.change(screen.getByLabelText(/evaluations/i), {
            target: { value: "1000" },
        });

        fireEvent.click(screen.getByRole("button", { name: /add variable/i }));

        const nameInputs = screen.getAllByLabelText(/^name$/i);
        const lowerInputs = screen.getAllByLabelText(/lower/i);
        const upperInputs = screen.getAllByLabelText(/upper/i);

        fireEvent.change(nameInputs[0], { target: { value: "x" } });
        fireEvent.change(lowerInputs[0], { target: { value: "0" } });
        fireEvent.change(upperInputs[0], { target: { value: "1" } });

        fireEvent.change(nameInputs[1], { target: { value: "y" } });
        fireEvent.change(lowerInputs[1], { target: { value: "2" } });
        fireEvent.change(upperInputs[1], { target: { value: "3" } });

        fireEvent.click(screen.getByRole("button", { name: /create job/i }));

        await screen.findByText(/job created: abc-123 \(queued\)/i);

        expect(jobsApi.createJob).toHaveBeenCalledWith({
            name: "test-job",
            integrand: "x+y",
            variables: [
                { name: "x", lower: "0", upper: "1" },
                { name: "y", lower: "2", upper: "3" },
            ],
            evaluations: 1000,
        });
    });

    it("sanitizes the job name input in the UI", () => {
        render(<JobForm />);

        const jobNameInput = screen.getByLabelText(/job name/i) as HTMLInputElement;

        fireEvent.change(jobNameInput, {
            target: { value: "job name!_test-1" },
        });

        expect(jobNameInput.value).toBe("jobname_test-1");
    });

    it("sanitizes the integrand input in the UI", () => {
        render(<JobForm />);

        const integrandInput = screen.getByLabelText(/integrand/i) as HTMLTextAreaElement;

        fireEvent.change(integrandInput, {
            target: { value: "sin(x); a@b#c" },
        });

        expect(integrandInput.value).toBe("sin(x) abc");
    });

    it("sanitizes the integration variable name input in the UI", () => {
        render(<JobForm />);

        const variableNameInput = screen.getByLabelText(/^name$/i) as HTMLInputElement;

        fireEvent.change(variableNameInput, {
            target: { value: "x-1!_abc" },
        });

        expect(variableNameInput.value).toBe("x1_abc");
    });

    it("limits the integration variable name length in the UI", () => {
        render(<JobForm />);

        const variableNameInput = screen.getByLabelText(/^name$/i) as HTMLInputElement;

        fireEvent.change(variableNameInput, {
            target: { value: "abcdefghijklmnopqr" },
        });

        expect(variableNameInput.value).toBe("abcdefghijklmnop");
    });

    it("sanitizes the lower boundary input in the UI", () => {
        render(<JobForm />);

        const lowerInput = screen.getByLabelText(/lower/i) as HTMLInputElement;

        fireEvent.change(lowerInput, {
            target: { value: "a1b2.3c" },
        });

        expect(lowerInput.value).toBe("12.3");
    });

    it("sanitizes the upper boundary input in the UI", () => {
        render(<JobForm />);

        const upperInput = screen.getByLabelText(/upper/i) as HTMLInputElement;

        fireEvent.change(upperInput, {
            target: { value: "x9y8.7z" },
        });

        expect(upperInput.value).toBe("98.7");
    });

    it("submits evaluations as a number", async () => {
        vi.mocked(jobsApi.createJob).mockResolvedValue({
            jobId: "abc-123",
            status: "queued",
        });

        render(<JobForm />);

        fireEvent.change(screen.getByLabelText(/job name/i), {
            target: { value: "test-job" },
        });

        fireEvent.change(screen.getByLabelText(/integrand/i), {
            target: { value: "x^2" },
        });

        fireEvent.change(screen.getByLabelText(/evaluations/i), {
            target: { value: "5000" },
        });

        fireEvent.click(screen.getByRole("button", { name: /create job/i }));

        await screen.findByText(/job created: abc-123 \(queued\)/i);

        expect(jobsApi.createJob).toHaveBeenCalledWith(
            expect.objectContaining({
                evaluations: 5000,
            })
        );
    });

    it("clears the old error message after a successful submit", async () => {
        vi.mocked(jobsApi.createJob)
            .mockRejectedValueOnce(new Error("request failed"))
            .mockResolvedValueOnce({
                jobId: "abc-123",
                status: "queued",
            });

        render(<JobForm />);

        fireEvent.change(screen.getByLabelText(/job name/i), {
            target: { value: "test-job" },
        });

        fireEvent.change(screen.getByLabelText(/integrand/i), {
            target: { value: "x^2" },
        });

        fireEvent.click(screen.getByRole("button", { name: /create job/i }));
        expect(await screen.findByText(/failed to create job/i)).toBeInTheDocument();

        fireEvent.click(screen.getByRole("button", { name: /create job/i }));
        expect(await screen.findByText(/job created: abc-123 \(queued\)/i)).toBeInTheDocument();
        expect(screen.queryByText(/failed to create job/i)).not.toBeInTheDocument();
    });

    it("clears the old success message after a failed submit", async () => {
        vi.mocked(jobsApi.createJob)
            .mockResolvedValueOnce({
                jobId: "abc-123",
                status: "queued",
            })
            .mockRejectedValueOnce(new Error("request failed"));

        render(<JobForm />);

        fireEvent.change(screen.getByLabelText(/job name/i), {
            target: { value: "test-job" },
        });

        fireEvent.change(screen.getByLabelText(/integrand/i), {
            target: { value: "x^2" },
        });

        fireEvent.click(screen.getByRole("button", { name: /create job/i }));
        expect(await screen.findByText(/job created: abc-123 \(queued\)/i)).toBeInTheDocument();

        fireEvent.click(screen.getByRole("button", { name: /create job/i }));
        expect(await screen.findByText(/failed to create job/i)).toBeInTheDocument();
        expect(screen.queryByText(/job created: abc-123 \(queued\)/i)).not.toBeInTheDocument();
    });
});

