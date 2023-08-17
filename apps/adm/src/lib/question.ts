export type AltType = 2 | 4;
export type QuestionStyle = "image" | "audio" | "video" | "text";
export type QuestionCorrect = "1" | "2" | "3" | "4";

export interface ValidateFormDataProps {
    question: string;
    questionStyle: QuestionStyle;
    answers: string[];
    correct: QuestionCorrect;
    file: string;
}

export function validateFormData({ question, questionStyle, answers, correct, file }: ValidateFormDataProps): string | null {
    if (question == "") {
        return "A pergunta não pode estar em branco!";
    }

    if (questionStyle != "audio" &&
        questionStyle != "image" &&
        questionStyle != "text" &&
        questionStyle != "video"
    ) {
        return "Somente os estilos ('audio', 'image', 'text' ou 'video') são aceitos!";
    }

    if (questionStyle != "text" && file == "") {
        return "Para os estilos ('audio', 'image', ou 'video'), uma arquivo precisa ser fornecido";
    }

    const altType: AltType = answers[2] == "" && answers[3] == "" ? 2 : 4;

    if (altType == 2) {
        if (answers[0] == "") {
            return "A resposta 1 precisa ser preenchida!";
        } else if (answers[1] == "") {
            return "A reposta 2 precisa ser preenchida!";
        }

        if (correct != "1" && correct != "2") {
            return "Como a resposta 3 e 4 não foram fornecidas, apenas opções (1 ou 2) podem estar corretas";
        }
    } else if (altType == 4) {
        if (answers[2] == "") {
            return "A resposta 3 precisa ser preenchida!";
        } else if (answers[3] == "") {
            return "A reposta 4 precisa ser preenchida!";
        }

        if (correct != "1" && correct != "2" && correct != "3" && correct != "4") {
            return "Somente opções (1, 2, 3 ou 4) podem estar corretas!";
        }
    }

    return null;
}
