import { Option, Result } from "ts-results";

export enum StorageError {
    SerializationError,
    DeserializationError,
    SetError,
    GetError,
    DeleteError
}

export function displayFriendlyErr(e: StorageError): string {
    return "Algo deu errado! Erro " + e.toString();
}

export interface Repository<T> {
    set(key: string, item: T): Promise<Result<void, StorageError>>;
    get(key: string): Promise<Result<Option<T>, StorageError>>;
    delete(key: string): Promise<Result<boolean, StorageError>>;
    getAll(): Promise<Result<T[], StorageError>>;
}

export interface PrimitiveRepository {
    set(item: string): Promise<Result<void, StorageError>>;
    get(): Promise<Result<Option<string>, StorageError>>;
    delete(): Promise<Result<boolean, StorageError>>;
}
