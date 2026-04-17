export type VariableSpec = {
    name: string;
    lower: string;
    upper: string;
};

export type CreateJobRequest = {
    name: string;
    integrand: string;
    variables: VariableSpec[];
    evaluations: number;
};

export type CreateJobResponse = {
    jobId: string;
    status: string;
};