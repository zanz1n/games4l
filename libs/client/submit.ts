/**
 * @license AGPL-3.0
 * Licensed under GNU Affero General Public License
 * @author Izan Rodrigues <izanrodrigues999@gmail.com>
 * Project: https://github.com/games4l/repo
 */

const API_GATEWAY_URI = "https://api.games4l.izan.com.br";
const LOCAL_STORAGE_DUMP_KEY = "GAMES4L_FETCH_FAILED_DUMP";

export interface Question {
    quest: string;
    a1: string;
    a2: string;
    a3?: string;
    a4?: string;
    ccr: number;
    audio?: string;
    type: "4Alt" | "2Alt";
    file?: string;
    style: "audio" | "video" | "image" | "text";
    x?: number;
    y?: number;
}

export interface QuestTrackData {
    idx: number;
    order: Array<number>;
    start_time: number;
}

export class LifecycleManager {
    private _cache = new Map<string, QuestTrackData>();

    constructor(private _pacientName: string) {}
    
    async _submitTelemetry(questionIdx: number, timeTaken: number, answeredOrder: number[]) {
        const payload = {
            done_at: Date.now() - timeTaken,
            complete_time: timeTaken,
            answereds: answeredOrder,
            quest_id: questionIdx,
            pacient_name: this._pacientName
        };

        try {
            const res = await fetch(`${API_GATEWAY_URI}/telemetry`, {
                body: JSON.stringify(payload),
                method: "POST"
            });

            if (!res.ok) {
                throw new Error();
            }
        // Catching in case of network unavailable kind of error
        } catch (e) {
            let lsItem = localStorage.getItem(LOCAL_STORAGE_DUMP_KEY);

            if (!lsItem) {
                lsItem = "[]";
            }

            let parsed: unknown[];

            try {
                parsed = JSON.parse(lsItem);
            } catch (e) {
                parsed = [];
            }

            if (parsed.length > 50) {
                parsed = [];
            }

            parsed.push(payload);

            localStorage.setItem(LOCAL_STORAGE_DUMP_KEY, JSON.stringify(parsed));
        }
    }

    startQuestTrack(questIdx: string | number): string {
        let idx = 0;
        if (typeof questIdx == "string") {
            idx = Number(questIdx.charAt(1));
        } else if (typeof questIdx == "number") {
            idx = questIdx;
        } else throw new Error("Invalid argument questIdx, expected number or string");

        const uuid = crypto.randomUUID();

        this._cache.set(uuid, { idx, order: [], start_time: Date.now() });
        return uuid;
    }

    async onQuestAnswred(tID: string, answerID: number, isCorrect: boolean) {
        const quest = this._cache.get(tID);

        if (!quest) {
            throw new Error("Value is not tracked or is moved");
        }

        quest.order.push(answerID);

        if (isCorrect) {
            const timeOffset = Date.now() - quest.start_time;

            await this._submitTelemetry(quest.idx, timeOffset, quest.order);
            this._cache.delete(tID);
            return;
        }

        this._cache.set(tID, quest);
    }
}

export default function lifecycleManager(pacientName: string) {
    return new LifecycleManager(pacientName);
}
