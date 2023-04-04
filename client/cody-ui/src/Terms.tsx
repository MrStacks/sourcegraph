import React from 'react'

export const Terms: React.FunctionComponent<{
    acceptTermsButton?: JSX.Element
}> = ({ acceptTermsButton }) => (
    <div className="non-transcript-container">
        <p className="terms-header-container">Notice and Usage Policies</p>
        <div className="terms-container">
            <p>
                Sourcegraph Cody is an AI coding assistant that finds, explains, and writes code using context from your
                codebase.
            </p>
            <p>
                Accuracy: Cody uses context from your codebase to improve the accuracy of its responses compared to
                other AI-based tools. However, Sourcegraph does not guarantee the accuracy of Cody's answers.
            </p>
            <p>
                Ownership: Sourcegraph makes no claims of ownership over the code generated by Cody, nor does
                Sourcegraph claim ownership of the user's existing code. The user retains ownership of their code and
                responsibility for ensuring their code complies with software licenses and copyright law. Cody may make
                use of language models trained on large datasets of publicly available code. It is the user's
                responsibility to verify any code snippets emitted by Cody comply with existing software licenses and
                copyright law.
            </p>
            <p>
                Acceptable use: You must follow the acceptable use policies of the following LLM providers:{' '}
                <a href="https://www.anthropic.com/aup">Anthropic Acceptable Use Policy</a>
            </p>
        </div>

        <p className="terms-header-container">FAQs</p>
        <div className="terms-container">
            <p className="question">
                Q: Will my queries or Cody's responses to my queries be used as training data for any machine learning
                models?
            </p>
            <p className="answer">A: No.</p>
            <p className="question">Q: Will my queries to Cody be shared with any third parties?</p>
            <p className="answer">
                A: Yes, we will send your queries to LLM providers we use for the sole purpose of providing you the
                service.
            </p>
            <p className="question">Q: Will the LLMs use my Cody Q&A to train their models?</p>
            <p className="answer">
                A: No, Sourcegraph has obtained commitments from our LLM providers that no Cody requests will be used to
                train LLM models.
            </p>
            <p>
                <a className="cta" href="https://docs.sourcegraph.com/cody">
                    Learn more about Cody security standards and data retention
                </a>
                .
            </p>
        </div>
        {acceptTermsButton}
    </div>
)
