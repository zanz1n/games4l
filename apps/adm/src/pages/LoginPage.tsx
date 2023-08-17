import { useState } from "react";
import styles from "../components/form/Form.module.css";
import Header from "../components/Header";
import { Form, InputLabel, SubmitButton, SwitchPages } from "../components/form/Form";
import globals from "../lib/globals";
import { useAuth } from "../lib/Auth";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

interface SignInDomData {
    username: {
        value: string;
    };
    password: {
        value: string;
    };
}

function validate(target: unknown): target is SignInDomData {
    if (target &&
        typeof target == "object" &&
        "username" in target &&
        target["username"] &&
        typeof target["username"] == "object" &&
        "value" in target["username"] &&
        target["username"]["value"] &&
        typeof target["username"]["value"] == "string" &&
        "password" in target &&
        target["password"] &&
        typeof target["password"] == "object" &&
        "value" in target["password"] &&
        target["password"]["value"] &&
        typeof target["password"]["value"] == "string") {
        return true;
    }
    return false;
}

export default function LoginPage() {
    const navigate = useNavigate();
    const { logIn, isLoggedIn } = useAuth();

    useEffect(() => {
        if (isLoggedIn()) {
            navigate("/");
        }
    }, []);

    const [error, setErrorRaw] = useState<string | null>(null);

    function setError(e: string | null) {
        setErrorRaw(e);
        if (e == null) setSendable(true);
        else setSendable(false);
    }

    const [sendable, setSendable] = useState<boolean>(false);

    async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
        const target: unknown = e.target;

        if (validate(target)) {
            const username = target.username.value;
            const password = target.password.value;

            try {
                await logIn(username, password);
                setError(null);
                navigate("/");
            } catch (e) {
                e instanceof Error ? setError(e.message) : setError("Algo deu errado");
            }
        } else {
            setError("O usuário e a senha precisam ser informados.");
        }
    }

    function onValueUpdate() {
        const username = document.getElementById("username") as HTMLInputElement;
        const password = document.getElementById("password") as HTMLInputElement;
        if (!username.value || username.value == "" || !password.value || password.value == "") {
            setSendable(false);
            return;
        }
        else setSendable(true);
    }

    return (
        <>
            <Header />
            <main className={styles.main}>
                <div className={styles.formContainer}>
                    <div className={styles.formTitle}>
                        <h1>Log In</h1>
                    </div>
                    <Form error={error} onSubmit={onSubmit}>
                        <InputLabel required
                            onChange={onValueUpdate}
                            identifier="username"
                            type="text">Username ou Email
                        </InputLabel>
                        <InputLabel required
                            onChange={onValueUpdate}
                            identifier="password"
                            type="password">Senha
                        </InputLabel>

                        <SubmitButton enabled={sendable}>Log In</SubmitButton>
                        <SwitchPages
                            plain="Algum Problema? Contacte no"
                            to={globals.TeamsUrl}>Teams
                        </SwitchPages>
                    </Form>
                </div>
            </main>
        </>
    );
}
