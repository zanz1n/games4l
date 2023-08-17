export interface FormProps {
    error?: string | null;
    onSubmit?: (e: React.FormEvent<HTMLFormElement>) => void | Promise<void>;
    children: React.ReactElement | React.ReactElement[];
    id?: string;
}

export type InputLabelProps = React.PropsWithChildren<{
    type: React.HTMLInputTypeAttribute
    identifier: string;
    required?: boolean;
    withError?: boolean;
    onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
}>

export interface SubmitButtonProps {
    enabled?: boolean;
    children: React.ReactElement | React.ReactElement[] | string;
}

export interface SwitchPagesProps {
    to: string;
    plain: string;
}
