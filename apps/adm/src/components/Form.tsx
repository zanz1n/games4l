import styles from "./Form.module.css";

export interface FormProps {
    error?: string | null;
    onSubmit?: (e: React.FormEvent<HTMLFormElement>) => void | Promise<void>;
    children: React.ReactElement | React.ReactElement[];
}

export default function Form({ error, onSubmit, children }: FormProps) {
    return(
        <form className={styles.form} onSubmit={(e) => {
            e.preventDefault();
            onSubmit?.(e);
        }}>
            <div className={`${styles.topError} ${error ? "" : styles.invisible}`}>
                <p>{error ?? "-"}</p>
            </div>
            {children}
            <div className={`${styles.topError} ${styles.invisible}`}>
                <p>-</p>
            </div>
        </form>
    );
}
