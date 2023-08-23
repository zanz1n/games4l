import DashMenuItem from "../../components/DashMenuItem";
import Header from "../../components/Header";
import styles from "./index.module.css";
import { useEffect } from "preact/hooks";
import CreateQuestionMenu from "./CreateQuestion";
import { AuthService } from "../../lib/Auth";
import { useRouter } from "preact-router";

export interface DashBoardProps {
    route: string;
}

export default function DashBoardMain() {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [route, navigate] = useRouter();

    const auth = AuthService.getInstance();

    useEffect(() => {
        if (!auth.isLoggedIn()) {
            navigate("/auth/login");
        }
    }, []);

    return <>
        <Header />
        <div className={styles.main}>
            <div className={styles.dashMenuContainer}>
                <div className={styles.menu}>
                    <DashMenuItem
                        to="/dash/create-question"
                        icon="/pencil.png"
                        text="Criar Pergunta"
                    />
                    <DashMenuItem
                        to="/dash/edit-question"
                        icon="/magnifier.png"
                        text="Editar Pergunta"
                    />
                </div>
            </div>
            <main className={styles.dashBody}>
                {
                    route.path == "/dash/create-question"
                        ?
                        <CreateQuestionMenu />
                        : ""
                }
            </main>
        </div>
    </>;
}
