export interface SubmitButtonProps {
    enabled?: boolean;
    children: React.ReactElement | React.ReactElement[] | string;
}

export default function SubmitButton({ children, enabled }: SubmitButtonProps) {
    return (
        <button disabled={!(enabled ?? true)} type="submit" >{children}</button>
    );
}
