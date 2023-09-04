import { Err, Ok, None } from "ts-results";
import { createCurrentSessionRepository, createSessionRepository } from "./impl";
import type { Option, Result } from "ts-results";
import { StorageError, type PrimitiveRepository, type Repository } from "./Repository";

export interface Session {
    name: string;
    createdAt: number;
    age: number;
    currentQuestion: number;
    hits: number;
    errors: number;
    skiped: number;
}

export interface SessionUpdateData {
    currentQuestion?: number;
    hits?: number;
    errors?: number;
    skiped?: number;
}

export class SessionManager {
    private static instance: SessionManager;

    public static getInstance(): SessionManager {
        if (!this.instance) {
            const currentSessionRepository = createCurrentSessionRepository();
            const sessionRepository = createSessionRepository();

            this.instance = new this(currentSessionRepository, sessionRepository);
        }

        return this.instance;
    }

    constructor(
        private currentSessionRepository: PrimitiveRepository,
        private sessionRepository: Repository<Session>,
    ) { }

    async getCurrent(): Promise<Result<Option<Session>, StorageError>> {
        const result = await this.currentSessionRepository.get();

        if (result.err) {
            return Err(result.val);
        }
        if (result.val.none) {
            return Ok(None);
        }

        const item = await this.sessionRepository.get(result.val.val);

        return item;
    }

    async setCurrent(name: string): Promise<Result<void, StorageError>> {
        const result = await this.currentSessionRepository.set(name);

        return result;
    }

    async create(name: string, age: number): Promise<Result<void, StorageError>> {
        const session = {
            name,
            age,
            createdAt: Date.now(),
            currentQuestion: 0,
            errors: 0,
            hits: 0,
            skiped: 0,
        } satisfies Session;

        const result = await this.sessionRepository.set(name, session);

        return result;
    }

    async update(name: string, data: SessionUpdateData): Promise<Result<Session, StorageError>> {
        const sess = await this.sessionRepository.get(name);
        if (sess.err) {
            return Err(sess.val);
        }
        if (sess.val.none) {
            return Err(StorageError.NotFound);
        }

        const { val } = sess.val;

        if (data.currentQuestion) {
            val.currentQuestion = data.currentQuestion;
        }
        if (data.errors) {
            val.errors = data.errors;
        }
        if (data.hits) {
            val.hits = data.hits;
        }
        if (data.skiped) {
            val.skiped = data.skiped;
        }

        const result = await this.sessionRepository.set(name, val);

        if (result.err) {
            return Err(result.val);
        }

        return Ok(val);
    }
}
