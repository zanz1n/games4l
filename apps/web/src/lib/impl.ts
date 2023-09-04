import { LocalStorage, LocalStoragePrimitive } from "./LocalStorageRepository";
import type { PrimitiveRepository, Repository } from "./Repository";
import type { Session } from "./Session";

export function createSessionRepository(): Repository<Session> {
    return new LocalStorage("sessions");
}

export function createCurrentSessionRepository(): PrimitiveRepository {
    return new LocalStoragePrimitive("current-session");
}
