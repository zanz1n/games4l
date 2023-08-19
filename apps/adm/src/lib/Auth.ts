import jwtDecode from "jwt-decode";
import globals from "./globals";

export class UnauthorizedErr extends Error {
    name = "UnauthorizedErr";
}

export interface UserJwtInfo {
    id: string;
    username: string;
    role: string;
}

export class AuthService {
    private static INSTANCE: AuthService;

    static getInstance(): AuthService {
        if (!this.INSTANCE) {
            this.INSTANCE = new this();
        }

        return this.INSTANCE;
    }

    private lsItem = "GAMES4L_TOKEN";

    getInfo() {
        const token = localStorage.getItem(this.lsItem);
        if (!token) {
            return null;
        }
        const data = jwtDecode(token);

        if (!isJwtPayload(data)) {
            return null;
        }
        return data;
    }

    isLoggedIn() {
        const item = localStorage.getItem(this.lsItem);

        if (!item) {
            return false;
        }

        if (item.length < 10) {
            localStorage.removeItem(this.lsItem);
            return false;
        }

        return true;
    }

    async logIn(credential: string, passwd: string) {
        try {
            const body = {
                credential,
                password: passwd,
            };

            const req = await fetch(globals.ApiGatewayUri + "/auth/signin", {
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
                localStorage.setItem(this.lsItem, data.token);
            } else throw new Error;

        } catch (e: unknown) {
            if (e instanceof Error) {
                throw new UnauthorizedErr(e.message);
            } else {
                throw new UnauthorizedErr("Autorização falhou");
            }
        }
    }

    logOut() {
        localStorage.removeItem(this.lsItem);
    }
}

function isJwtPayload(payload: unknown): payload is UserJwtInfo {
    return (typeof payload == "object" && payload &&
        "id" in payload && typeof payload["id"] &&
        "username" in payload && typeof payload["username"] == "string" &&
        "role" in payload && typeof payload["role"] == "string") == true;
}
