import { describe, expect, it } from "vitest";
import {
    sanitizeBoundary,
    sanitizeIntegrand,
    sanitizeIntegrationVariableName,
    sanitizeJobName,
} from "./inputSanitizers";

describe("sanitizeIntegrand", () => {
    it("keeps allowed characters", () => {
        expect(sanitizeIntegrand("(1+x^2+y^2+Pi*ln(1+z^2+2*x*y))^4")).toBe(
            "(1+x^2+y^2+Pi*ln(1+z^2+2*x*y))^4"
        );
    });

    it("removes disallowed characters", () => {
        expect(sanitizeIntegrand("sin(x); DROP TABLE")).toBe("sin(x) DROP TABLE");
    });
});

describe("sanitizeIntegrationVariableName", () => {
    it("keeps letters digits and underscore", () => {
        expect(sanitizeIntegrationVariableName("x_12")).toBe("x_12");
    });

    it("removes disallowed characters", () => {
        expect(sanitizeIntegrationVariableName("x-1!")).toBe("x1");
    });

    it("limits length to 16", () => {
        expect(sanitizeIntegrationVariableName("abcdefghijklmnopqr")).toBe(
            "abcdefghijklmnop"
        );
    });
});

describe("sanitizeBoundary", () => {
    it("keeps digits and decimal points", () => {
        expect(sanitizeBoundary("12.34")).toBe("12.34");
    });

    it("removes other characters", () => {
        expect(sanitizeBoundary("a1b2.3c")).toBe("12.3");
    });
});

describe("sanitizeJobName", () => {
    it("keeps letters digits underscore and hyphen", () => {
        expect(sanitizeJobName("job_1-test")).toBe("job_1-test");
    });

    it("removes disallowed characters", () => {
        expect(sanitizeJobName("job name!")).toBe("jobname");
    });
});