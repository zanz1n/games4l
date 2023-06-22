import { useNavigate } from "react-router-dom";
import styles from "../pages/dashboard/index.module.css";

export interface DashMenuItemProps {
    text: string;
    icon: string;
    to: string;
}

export default function DashMenuItem({ icon, text, to }: DashMenuItemProps) {
    const navigate = useNavigate();

    return <button className={styles.item} onClick={() => {navigate(to);}}>
        <img width="32px" height="32px" src={icon} alt={text+"-IMG"} />
        <div className={styles.itemText}>
            <p>{text}</p>
        </div>
    </button>;
}
