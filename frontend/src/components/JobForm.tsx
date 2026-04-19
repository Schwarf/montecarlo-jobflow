import React, { useState } from "react";
import { createJob } from "../api/jobs";
import type { VariableSpec } from "../types/job";
import {
    sanitizeBoundary,
    sanitizeIntegrand,
    sanitizeIntegrationVariableName,
} from "../utils/inputSanitizers";

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

    function handleRemoveIntegrationVariable(indexToRemove: number) {
        setIntegrationVariables(
            integrationVariables.filter((_, index) => index !== indexToRemove)
        );
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

            setResultMessage(`Job created: ${response.jobId} (${response.status})`);
        } catch (error) {
            setErrorMessage("Failed to create job");
            console.error(error);
        }
    }

    return (
        <form className="job-form" onSubmit={handleSubmit}>
            <div className="job-form-left">
                <section className="panel">
                    <h2> Job configuration</h2>

                    <div className="form-field">
                        <label htmlFor="name"> Job Name</label>
                        <input id="name" value={name} onChange={(event) => setName(event.target.value)}
                        />
                    </div>
                    <div className="form-field">
                        <label htmlFor="evaluations">Evaluations</label>
                        <input
                            id="evaluations"
                            type="number"
                            value={evaluations}
                            onChange={(event) => setEvaluations(Number(event.target.value))}
                        />
                    </div>

                    <div className="help-box">
                        <h3>Allowed job name format</h3>
                        <p>
                            Keep the name short and descriptive.
                        </p>
                    </div>
                </section>

                <section className="panel">
                    <h2>Integration variables</h2>

                    <div className="help-box">
                        <h3>Allowed variable format</h3>
                        <p>
                            Variable names may contain letters, digits, and underscores. Maximum length
                            is 16 characters.
                        </p>
                        <p>
                            Bounds currently only allow digits and decimal points.
                        </p>
                    </div>

                    {integrationVariables.map((variable, index) => (
                        <div className="variable-card" key={index}>
                            <div className="variable-row">
                                <div className="form-field">
                                    <label htmlFor={`var-name-${index}`}>Name</label>
                                    <input
                                        id={`var-name-${index}`}
                                        value={variable.name}
                                        maxLength={16}
                                        onChange={(event) =>
                                            handleIntegrationVariableChange(
                                                index,
                                                "name",
                                                sanitizeIntegrationVariableName(event.target.value)
                                            )
                                        }
                                    />
                                </div>

                                <div className="form-field">
                                    <label htmlFor={`var-lower-${index}`}>Lower</label>
                                    <input
                                        id={`var-lower-${index}`}
                                        value={variable.lower}
                                        onChange={(event) =>
                                            handleIntegrationVariableChange(
                                                index,
                                                "lower",
                                                sanitizeBoundary(event.target.value)
                                            )
                                        }
                                    />
                                </div>

                                <div className="form-field">
                                    <label htmlFor={`var-upper-${index}`}>Upper</label>
                                    <input
                                        id={`var-upper-${index}`}
                                        value={variable.upper}
                                        onChange={(event) =>
                                            handleIntegrationVariableChange(
                                                index,
                                                "upper",
                                                sanitizeBoundary(event.target.value)
                                            )
                                        }
                                    />
                                </div>
                            </div>

                            <button
                                type="button"
                                onClick={() => handleRemoveIntegrationVariable(index)}
                                disabled={integrationVariables.length === 1}
                            >
                                Remove this variable
                            </button>
                        </div>
                    ))}

                    <button type="button" onClick={handleAddIntegrationVariable}>
                        Add variable
                    </button>
                </section>
                <section className="panel">
                    <button type="submit">Create job</button>

                    {resultMessage && <p>{resultMessage}</p>}
                    {errorMessage && <p>{errorMessage}</p>}
                </section>
            </div>

            <div className="job-form-right">
                <section className="panel integrand-panel">
                    <h2>Integrand</h2>

                    <textarea
                        id="integrand"
                        className="integrand-textarea"
                        value={integrand}
                        onChange={(event) =>
                            setIntegrand(sanitizeIntegrand(event.target.value))
                        }
                    />

                    <div className="help-box">
                        <h3>Supported integrand format</h3>
                        <p>
                            Allowed characters: letters, digits, parentheses, plus, minus,
                            multiply, divide, power, decimal point, and spaces.
                        </p>
                        <p><strong>Supported functions</strong></p>
                        <p>Trigonometric: sin, cos, tan</p>
                        <p>Inverse trigonometric: asin, acos, atan</p>
                        <p>Hyperbolic: sinh, cosh, tanh</p>
                        <p>Inverse hyperbolic: asinh, acosh, atanh</p>
                        <p>Logarithmic and exponential: log10, ln, log2, exp</p>

                        <p><strong>Supported constants</strong></p>
                        <p>Circular constant: Pi</p>
                        <p>Euler's number: E</p>
                    </div>
                </section>
            </div>
        </form>
    );
}