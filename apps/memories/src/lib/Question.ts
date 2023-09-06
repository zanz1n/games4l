export interface Question {
    question: string;
    answers: string[];
    correct_answer: number;
    type: "2Alt" | "4Alt";
    style: "image" | "audio" | "video" | "text";
    difficulty: number;
    file: string | null;
    image_width: number | null;
    image_height: number | null;
}

async function getJson() {
    const res = await fetch("/questions.json");

    if (!res.ok) {
        throw new Error("HTTP code " + res.statusText);
    }

    const json = await res.json();

    return json;
}

const json = (await getJson()) as Question[];

export function questionsLength(): number {
    return json.length;
}

export function getQuestion(id: number): Question | null {
    return json[id] as Question;
}
