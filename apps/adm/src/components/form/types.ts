import { ComponentChildren, JSX } from "preact";

export interface FormProps {
    error?: string | null;
    onSubmit?: JSX.GenericEventHandler<HTMLFormElement>;
    children: ComponentChildren;
    id?: string;
}

export interface InputLabelProps {
    type: string;
    identifier: string;
    required?: boolean;
    withError?: boolean;
    onChange?: JSX.GenericEventHandler<HTMLInputElement>;
    children: ComponentChildren;
}

export interface SubmitButtonProps {
    enabled?: boolean;
    children: ComponentChildren;
}

export interface SwitchPagesProps {
    to: string;
    plain: string;
}
