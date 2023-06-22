import { Route, createBrowserRouter, createRoutesFromElements, useNavigate } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import App from "./App";
import DashBoardMain from "./pages/dashboard";
import { useEffect } from "react";

function RedirectTo({ to }: { to: string }) {
    const navigate = useNavigate();

    useEffect(() => {
        navigate(to);
    }, []);

    return <h1 style={{
        textAlign: "center",
        marginTop: "32px"
    }}>Redirecting...</h1>;
}

export const router = createBrowserRouter(createRoutesFromElements(
    <Route path="/" element={<App/>}>
        <Route path="/" element={<RedirectTo to="/dash"/>} />
        <Route path="/auth/login" element={<LoginPage/>}/>
        <Route path="/dash" element={<DashBoardMain route="/" />}/>
        <Route
            path="/dash/create-question"
            element={<DashBoardMain route="create-question"/>}
        />
        <Route
            path="/dash/update-question"
            element={<DashBoardMain route="update-question"/>}
        />
    </Route>
));
