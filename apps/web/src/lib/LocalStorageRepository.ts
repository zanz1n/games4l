import { Ok, Err, None, Some } from "ts-results";
import { StorageError } from "./Repository";
import type { Result, Option } from "ts-results";
import type { PrimitiveRepository, Repository } from "./Repository";

export interface MapObject<T> {
    [key: string]: T | null;
}

export interface ToString {
    toString(): string;
}

export class LocalStorage<T> implements Repository<T> {
    constructor(private itemIdx: string) { }

    async set(key: string, item: T): Promise<Result<void, StorageError>> {
        try {
            let all = localStorage.getItem(this.itemIdx) ?? "{}";
            const arr = JSON.parse(all) as MapObject<T>;

            arr[key] = item;

            try {
                all = JSON.stringify(arr);
            } catch (_) {
                return Err(StorageError.SerializationError);
            }

            localStorage.setItem(this.itemIdx, all);

            return Ok(void 0);
        } catch {
            return Err(StorageError.SetError);
        }
    }

    async get(key: string): Promise<Result<Option<T>, StorageError>> {
        try {
            const all = localStorage.getItem(this.itemIdx) ?? "{}";
            const arr = JSON.parse(all) as MapObject<T>;

            const item = arr[key];

            if (!item) {
                return Ok(None);
            }

            return Ok(Some(item));

        } catch (_) {
            return Err(StorageError.GetError);
        }
    }

    async delete(key: string): Promise<Result<boolean, StorageError>> {
        try {
            let all = localStorage.getItem(this.itemIdx) ?? "{}";
            const arr = JSON.parse(all) as MapObject<T>;

            const has = !!arr[key];
            delete arr[key];

            try {
                all = JSON.stringify(arr);
            } catch (_) {
                return Err(StorageError.SerializationError);
            }
            localStorage.setItem(this.itemIdx, all);

            return Ok(has);
        } catch (_) {
            return Err(StorageError.DeleteError);
        }
    }

    async getAll(): Promise<Result<T[], StorageError>> {
        try {
            const all = localStorage.getItem(this.itemIdx) ?? "{}";
            const arr = JSON.parse(all) as MapObject<T>;

            const newArr = [] as T[];

            const keys = Object.keys(arr);
            for (const key of keys) {
                const item = arr[key];

                if (item) {
                    newArr.push(item);
                }
            }

            return Ok(newArr);
        } catch (_) {
            return Err(StorageError.GetError);
        }
    }
}

export class LocalStoragePrimitive implements PrimitiveRepository {
    constructor(private itemIdx: string) { }

    async set(item: string): Promise<Result<void, StorageError>> {
        try {
            localStorage.setItem(this.itemIdx, item);
            return Ok(void 0);
        } catch (_) {
            return Err(StorageError.SetError);
        }
    }

    async get(): Promise<Result<Option<string>, StorageError>> {
        try {
            const item = localStorage.getItem(this.itemIdx);

            if (!item) {
                return Ok(None);
            }

            return Ok(Some(item));
        } catch (_) {
            return Err(StorageError.SetError);
        }
    }

    async delete(): Promise<Result<boolean, StorageError>> {
        try {
            const has = !!localStorage.getItem(this.itemIdx);

            if (has) {
                localStorage.removeItem(this.itemIdx);
            }
            return Ok(has);
        } catch (_) {
            return Err(StorageError.DeleteError);
        }
    }
}
