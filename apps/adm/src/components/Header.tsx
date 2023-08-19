import { Link, useRouter } from "preact-router";
import styles from "./Header.module.css";
import { AuthService } from "../lib/Auth";

export default function Header() {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [_, navigate] = useRouter();

    const auth = AuthService.getInstance();

    const user = auth.getInfo();

    return (
        <div className={styles.header}>
            <header className={styles.headerContainer}>
                <div className={styles.left}>
                    <Link href="/">
                        <h1 className={styles.title}>Games4Life</h1>
                    </Link>
                </div>
                <div className={styles.right}>
                    {auth.isLoggedIn() ? (
                        <>
                            <nav className={styles.nav}>
                                <Link href="/">Home</Link>
                            </nav>
                            <button onClick={() => {
                                auth.logOut();
                                navigate("/auth/login");
                            }} className={styles.purple}>Logout</button>
                            <div className={styles.sideUsername}>
                                <p>Logged in as</p>
                                <p>{user?.username}</p>
                            </div>
                        </>
                    ) : (
                        <>
                            <button
                                onClick={() => navigate("/auth/login")}
                                className={styles.green}>Sign In
                            </button>
                        </>
                    )}
                </div>
            </header>
        </div>
    );
}
