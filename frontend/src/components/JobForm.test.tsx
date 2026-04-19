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
});
