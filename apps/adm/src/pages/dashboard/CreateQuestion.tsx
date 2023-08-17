import { useEffect, useState } from "react";
import { BigForm, InputLabel, SubmitButton } from "../../components/form/BigForm";
import styles from "../../components/form/BigForm.module.css";

type AltType = 2 | 4;
type QuestionStyle = "image" | "audio" | "video" | "text";

export default function CreateQuestionMenu() {
    const [error, setError] = useState<string | null>(null);

    const [question, setQuestion] = useState("");
    const [answers, setAnswers] = useState(["", "", "", ""]);
    const [correct, setCorrect] = useState("");
    const [questionStyle, setQuestionStyle] = useState("" as QuestionStyle);
    const [file, setFile] = useState("");

    useEffect(() => {
        console.log("updatedAt: " + Date.now().toString());

        if (question == "") {
            return setError("A pergunta não pode estar em branco!");
        }

        if (questionStyle != "audio" &&
            questionStyle != "image" &&
            questionStyle != "text" &&
            questionStyle != "video"
        ) {
            return setError("Somente os estilos ('audio', 'image', 'text' ou 'video') são aceitos!");
        }

        if (questionStyle != "text") {
            return setError("Para os estilos ('audio', 'image', ou 'video'), uma arquivo precisa ser fornecido");
        }

        const altType: AltType = answers[2] == "" && answers[3] == "" ? 2 : 4;

        if (altType == 2) {
            if (answers[0] == "") {
                return setError("A resposta 1 precisa ser preenchida!");
            } else if (answers[1] == "") {
                return setError("A reposta 2 precisa ser preenchida!");
            }

            if (correct != "1" && correct != "2") {
                return setError("Como a resposta 3 e 4 não foram fornecidas, apenas opções (1 ou 2) podem estar corretas");
            }
        } else if (altType == 4) {
            if (answers[2] == "") {
                return setError("A resposta 3 precisa ser preenchida!");
            } else if (answers[3] == "") {
                return setError("A reposta 4 precisa ser preenchida!");
            }

            if (correct != "1" && correct != "2" && correct != "3" && correct != "4") {
                return setError("Somente opções (1, 2, 3 ou 4) podem estar corretas!");
            }
        }

        setError(null);
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
                        onChange={(e) => { setCorrect(e.target.value); }}
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
