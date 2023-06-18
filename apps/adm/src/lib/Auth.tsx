import { createContext, useMemo } from "react";
import globals from "./globals";
import jwtDecode from "jwt-decode";
import { useContext } from "react";

export class UnauthorizedErr extends Error {
    name = "UnauthorizedErr";
}

const localStorageItemName = "GAMES4L_TOKEN";

export interface UserJwtInfo {
    id: string;
    username: string;
    role: string;
}

export interface AuthContext {
    isLoggedIn(): boolean
    logIn(credential: string, passwd: string): Promise<void>
    logOut(): void
    getInfo(): UserJwtInfo | null
}

const Context = createContext({} as AuthContext);

function isJwtPayload(payload: unknown): payload is UserJwtInfo {
    return (typeof payload == "object" && payload &&
    "id" in payload && typeof payload["id"] &&
    "username" in payload && typeof payload["username"] == "string" &&
    "role" in payload && typeof payload["role"] == "string") == true;
}

export default function AuthProvider({children}: { children: any }) {
    function getInfo(): UserJwtInfo | null {
        const token = localStorage.getItem(localStorageItemName);
        if (!token) {
            return null;
        }
        const data = jwtDecode(token);

        if (!isJwtPayload(data)) {
            return null;
        }
        return data;
    }

    function isLoggedIn() {
        const item = localStorage.getItem(localStorageItemName);

        if (!item) {
            return false;
        }

        if (item.length < 10) {
            localStorage.removeItem(localStorageItemName);
            return false;
        }

        return true;
    }

    async function logIn(credential: string, passwd: string) {
        try {
            const body = {
                credential,
                password: passwd,
            };

            const req = await fetch(globals.ApiGatewayUri+"/auth/signin", {
                method: "POST",
                body: JSON.stringify(body),
                headers: {
                    "Content-Type": "application/json",
                    "Accept": "application/json"
                }
            });

            if (!req.ok) {
                throw new Error;
            }

            const data: unknown = await req.json();

            if (data && typeof data == "object" && "token" in data && typeof data["token"] == "string") {
                localStorage.setItem(localStorageItemName, data.token);
            } else throw new Error;

        } catch(e: any) {
            throw new UnauthorizedErr(e["message"] ?? "Autorização falhou");
        }
    }

    function logOut() {
        localStorage.removeItem(localStorageItemName);
    }

    const aCtxValue = useMemo(() => ({
        getInfo,
        isLoggedIn,
        logIn,
        logOut
    } satisfies AuthContext), []);

    return <Context.Provider value={aCtxValue}>
        {children}
    </Context.Provider>;
}

export function useAuth() {
    return useContext(Context);
}
