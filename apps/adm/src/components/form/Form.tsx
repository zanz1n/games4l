import { Link } from "preact-router";
import { PropsWithChildren } from "preact/compat";
import styles from "./Form.module.css";
import { FormProps, InputLabelProps, SubmitButtonProps, SwitchPagesProps } from "./types";

export function SwitchPages({ plain, children, to }: PropsWithChildren<SwitchPagesProps>) {
    return (
        <div className={styles.switchPages}>
            <p>{plain}<> </><Link href={to}>{children}</Link></p>
        </div>
    );
}

export function SubmitButton({ children, enabled }: SubmitButtonProps) {
    return (
        <button disabled={!(enabled ?? true)} type="submit" >{children}</button>
    );
}

export function InputLabel({ required, type, identifier, children, onChange }: InputLabelProps) {
    return (
        <div className={styles.inputLabel}>
            <label htmlFor={identifier}>{children}</label>
            <div className={styles.formInput}>
                <input onInput={onChange} required={required} type={type} name={identifier} id={identifier} />
            </div>
        </div>
    );
}

export function Form({ error, onSubmit, children, id }: FormProps) {
    return (
        <form id={id} className={styles.form} onSubmit={(e) => {
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
