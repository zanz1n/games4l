import styles from "./BigForm.module.css";
import otherStyles from "./Form.module.css";
import { FormProps } from "./Form";
import { InputLabelProps } from "./InputLabel";
import { SubmitButtonProps } from "./SubmitButton";

export function SubmitButton({ children, enabled }: SubmitButtonProps) {
    return (
        <button className={styles.formButton} disabled={!(enabled ?? true)} type="submit" >{children}</button>
    );
}

export function InputLabel({ identifier, onChange, children, required, type }: InputLabelProps) {
    return (
        <div className={styles.inputLabel}>
            <label htmlFor={identifier}>{children}</label>
            <div className={styles.formInput}>
                <input onChange={onChange} required={required} type={type} name={identifier} id={identifier} />
            </div>
        </div>
    );
}

export function BigForm({ error, onSubmit, children }: FormProps) {
    return(
        <form className={styles.form} onSubmit={(e) => {
            e.preventDefault();
            onSubmit?.(e);
        }}>
            <div className={`${otherStyles.topError} ${error ? "" : otherStyles.invisible}`}>
                <p>{error ?? "-"}</p>
            </div>
            {children}
            <div className={`${otherStyles.topError} ${otherStyles.invisible}`}>
                <p>-</p>
            </div>
        </form>
    );
}
