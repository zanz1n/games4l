import { Link, useNavigate } from "react-router-dom";
import styles from "./Header.module.css";
import { useAuth } from "../lib/Auth";

export default function Header() {
    const { isLoggedIn, getInfo, logOut } = useAuth();

    const user = getInfo();

    const navigate = useNavigate();

    return (
        <div className={styles.header}>
            <header className={styles.headerContainer}>
                <div className={styles.left}>
                    <Link to="/">
                        <h1 className={styles.title}>Games4Life</h1>
                    </Link>
                </div>
                <div className={styles.right}>
                    {
                        isLoggedIn() ? (
                            <>
                                <nav className={styles.nav}>
                                    <Link to="/">Home</Link>
                                    <Link to="/dash/questions">Perguntas</Link>
                                    <Link to="/dash/users">Usuários</Link>
                                </nav>
                                <button onClick={() => {
                                    logOut();
                                    navigate("/auth/login");
                                }} className={styles.purple}>Logout</button>
                                <div className={styles.sideUsername}>
                                    <p>Logged in as</p>
                                    <p>{user?.username}</p>
                                </div>
                            </>
                        ) : (
                            <>
                                <button onClick={() => navigate("/auth/login")} className={styles.green}>Sign In</button>
                            </>
                        )
                    }
                </div>
            </header>
        </div>
    );
}
