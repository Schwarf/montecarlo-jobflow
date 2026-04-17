import React, { useState } from "react";
import { createJob } from "../api/jobs";
import type { VariableSpec } from "../types/job";

export function JobForm() {
    const [name, setName] = useState("");
    const [integrand, setIntegrand] = useState("");
    const [evaluations, setEvaluations] = useState(1000000);
    const [integrationVariables, setIntegrationVariables] = useState<VariableSpec[]>([
        { name: "x", lower: "0", upper: "1" },
    ]);
    const [resultMessage, setResultMessage] = useState("");
    const [errorMessage, setErrorMessage] = useState("");

    function handleIntegrationVariableChange(
        index: number,
        field: keyof VariableSpec,
        value: string
    ) {
        const updatedVariables = [...integrationVariables];
        updatedVariables[index] = {
            ...updatedVariables[index],
            [field]: value,
        };
        setIntegrationVariables(updatedVariables);
    }

    function handleAddIntegrationVariable() {
        setIntegrationVariables([
            ...integrationVariables,
            { name: "", lower: "", upper: "" },
        ]);
    }

    async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();
        setResultMessage("");
        setErrorMessage("");

        try {
            const response = await createJob({
                name,
                integrand,
                variables: integrationVariables,
                evaluations,
            });

            setResultMessage(`Job created: ${response.jobID} (${response.status})`);
        } catch (error) {
            setErrorMessage("Failed to create job");
            console.error(error);
        }
    }

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="name">Job name</label>
                <input
                    id="name"
                    value={name}
                    onChange={(event) => setName(event.target.value)}
                />
            </div>

            <div>
                <label htmlFor="integrand">Integrand</label>
                <input
                    id="integrand"
                    value={integrand}
                    onChange={(event) => setIntegrand(event.target.value)}
                />
            </div>

            <div>
                <label htmlFor="evaluations">Evaluations</label>
                <input
                    id="evaluations"
                    type="number"
                    value={evaluations}
                    onChange={(event) => setEvaluations(Number(event.target.value))}
                />
            </div>

            <h2>Integration variables</h2>

            {integrationVariables.map((variable, index) => (
                <div key={index}>
                    <div>
                        <label htmlFor={`var-name-${index}`}>Name</label>
                        <input
                            id={`var-name-${index}`}
                            value={variable.name}
                            onChange={(event) =>
                                handleIntegrationVariableChange(index, "name", event.target.value)
                            }
                        />
                    </div>

                    <div>
                        <label htmlFor={`var-lower-${index}`}>Lower bound</label>
                        <input
                            id={`var-lower-${index}`}
                            value={variable.lower}
                            onChange={(event) =>
                                handleIntegrationVariableChange(index, "lower", event.target.value)
                            }
                        />
                    </div>

                    <div>
                        <label htmlFor={`var-upper-${index}`}>Upper bound</label>
                        <input
                            id={`var-upper-${index}`}
                            value={variable.upper}
                            onChange={(event) =>
                                handleIntegrationVariableChange(index, "upper", event.target.value)
                            }
                        />
                    </div>
                </div>
            ))}

            <button type="button" onClick={handleAddIntegrationVariable}>
                Add variable
            </button>

            <div>
                <button type="submit">Create job</button>
            </div>

            {resultMessage && <p>{resultMessage}</p>}
            {errorMessage && <p>{errorMessage}</p>}
        </form>
    );
}
