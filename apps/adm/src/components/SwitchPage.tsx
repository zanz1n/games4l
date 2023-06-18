import styles from "./Form.module.css";
import { Link } from "react-router-dom";

export interface SwitchPagesProps {
    to: string;
    plain: string;
}

export default function SwitchPages({ plain, children, to }: React.PropsWithChildren<SwitchPagesProps>) {
    return (
        <div className={styles.switchPages}>
            <p>{plain}<> </><Link to={to}>{children}</Link></p>
        </div>
    );
}
