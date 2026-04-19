import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { JobForm } from "./JobForm";

describe("JobForm", () => {
    it("renders the main form fields", () => {
        render(<JobForm />);

        expect(screen.getByLabelText(/job name/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/evaluations/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/integrand/i)).toBeInTheDocument();
        expect(screen.getByRole("button", { name: /create job/i })).toBeInTheDocument();
    });
});