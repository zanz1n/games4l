import { Ok, Err, None, Some } from "ts-results";
import { StorageError } from "./Repository";
import type { Result, Option } from "ts-results";
import type { Repository } from "./Repository";

export interface MapObject<T> {
    [key: string]: T | null;
}

export class LocalStorageRepository<T> implements Repository<T> {
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
