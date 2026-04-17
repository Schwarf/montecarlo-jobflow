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
    jobID: string;
    status: string;
};