import { useNavigate } from "react-router-dom";
import DashMenuItem from "../../components/DashMenuItem";
import Header from "../../components/Header";
import styles from "./index.module.css";
import { useAuth } from "../../lib/Auth";
import { useEffect } from "react";
import CreateQuestionMenu from "./CreateQuestion";

export interface DashBoardProps {
    route: string;
}

export default function DashBoardMain({ route }: DashBoardProps) {
    const { isLoggedIn } = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        if (!isLoggedIn()) {
            navigate("/auth/login");
        }
    }, []);

    return <>
        <Header/>
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
                    route == "create-question"
                        ?
                        <CreateQuestionMenu/>
                        : ""
                }
            </main>
        </div>
    </>;
}
