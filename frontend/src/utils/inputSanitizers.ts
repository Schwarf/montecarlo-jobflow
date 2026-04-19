export function sanitizeIntegrand(value: string): string {
    return value.replace(/[^a-zA-Z0-9()+\-*/^.,\s]/g, "");
}

export function sanitizeIntegrationVariableName(value: string): string {
    return value.replace(/[^a-zA-Z0-9_]/g, "").slice(0, 16);
}

export function sanitizeBoundary(value: string): string {
    return value.replace(/[^0-9.]/g, "");
}

export function sanitizeJobName(value: string): string {
    return value.replace(/[^a-zA-Z0-9_-]/g, "");
}