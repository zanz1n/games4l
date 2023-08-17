import { useEffect, useState } from "react";
import { BigForm, InputLabel, SubmitButton } from "../../components/form/BigForm";
import styles from "../../components/form/BigForm.module.css";
import { QuestionCorrect, QuestionStyle, validateFormData } from "../../lib/question";

export default function CreateQuestionMenu() {
    const [error, setError] = useState<string | null>(null);

    const [question, setQuestion] = useState("");
    const [answers, setAnswers] = useState(["", "", "", ""]);
    const [correct, setCorrect] = useState("" as QuestionCorrect);
    const [questionStyle, setQuestionStyle] = useState("" as QuestionStyle);
    const [file, setFile] = useState("");

    useEffect(() => {
        setError(validateFormData({ answers, correct, question, questionStyle, file }));
    }, [question, answers, correct, questionStyle, file]);

    function updateAnswer(idx: number) {
        return (e: React.ChangeEvent<HTMLInputElement>) => {
            const mutation = [...answers];
            mutation[idx] = e.target.value;

            setAnswers(mutation);
        };
    }

    return <div className={styles.formBigContainer}>
        <div>
            <h1>Criar pergunta</h1>
        </div>
        <BigForm error={error} id="create-question-form">
            <div className={styles.bigFormBody}>
                <div>
                    <InputLabel identifier="question" type="text"
                        onChange={(e) => { setQuestion(e.target.value); }}
                    >Pergunta</InputLabel>

                    <InputLabel identifier="correct_answer" type="text"
                        onChange={(e) => { setCorrect(e.target.value as QuestionCorrect); }}
                    >Resposta correta</InputLabel>

                    <InputLabel identifier="question_style" type="text"
                        onChange={(e) => { setQuestionStyle(e.target.value as QuestionStyle); }}
                    >Estilo</InputLabel>

                    <InputLabel identifier="file" type="text"
                        onChange={(e) => { setFile(e.target.value); }}
                    >Arquivo (imagem/video/audio)</InputLabel>
                </div>
                <div>
                    <InputLabel identifier="answer1" type="text"
                        onChange={updateAnswer(0)}
                    >Resposta 1</InputLabel>

                    <InputLabel identifier="answer2" type="text"
                        onChange={updateAnswer(1)}
                    >Resposta 2</InputLabel>

                    <InputLabel identifier="answer3" type="text"
                        onChange={updateAnswer(2)}
                    >Resposta 3 (opcional)</InputLabel>

                    <InputLabel identifier="answer4" type="text"
                        onChange={updateAnswer(3)}
                    >Resposta 4 (opcional)</InputLabel>
                </div>
            </div>
        </BigForm>
        <SubmitButton enabled={!error}>Publicar</SubmitButton>
    </div>;
}
