import { describe, expect, it, vi, beforeEach} from "vitest";
import {fireEvent, render, screen} from "@testing-library/react";
import { JobForm } from "./JobForm";
import * as jobsApi from "../api/jobs";


vi.mock("../api/jobs", () => ({
    createJob: vi.fn(),
}));

describe("JobForm", () => {

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders the main form fields", () => {
        render(<JobForm />);

        expect(screen.getByLabelText(/job name/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/evaluations/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/integrand/i)).toBeInTheDocument();
        expect(screen.getByRole("button", { name: /create job/i })).toBeInTheDocument();
    });

    it("adds another integration variable", () => {
        render(<JobForm />);

        expect(screen.getAllByLabelText(/^name$/i)).toHaveLength(1);

        fireEvent.click(screen.getByRole("button", { name: /add variable/i }));

        expect(screen.getAllByLabelText(/^name$/i)).toHaveLength(2);
    });

    it("removes an integration variable", () => {
        render(<JobForm />);

        fireEvent.click(screen.getByRole("button", { name: /add variable/i }));
        expect(screen.getAllByLabelText(/^name$/i)).toHaveLength(2);

        const removeButtons = screen.getAllByRole("button", { name: /remove this variable/i });
        fireEvent.click(removeButtons[1]);

        expect(screen.getAllByLabelText(/^name$/i)).toHaveLength(1);
    });

    it("submits the form and shows the success message", async () => {
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

        fireEvent.click(screen.getByRole("button", { name: /create job/i }));

        expect(await screen.findByText(/job created: abc-123 \(queued\)/i)).toBeInTheDocument();
        expect(jobsApi.createJob).toHaveBeenCalled();
    });

    it("shows an error message when submission fails", async () => {
        vi.mocked(jobsApi.createJob).mockRejectedValue(new Error("request failed"));

        render(<JobForm />);

        fireEvent.change(screen.getByLabelText(/job name/i), {
            target: { value: "test-job" },
        });

        fireEvent.change(screen.getByLabelText(/integrand/i), {
            target: { value: "x^2" },
        });

        fireEvent.click(screen.getByRole("button", { name: /create job/i }));

        expect(await screen.findByText(/Failed to create job/i)).toBeInTheDocument();
    });
});
