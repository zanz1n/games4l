import { Route, createBrowserRouter, createRoutesFromElements } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import App from "./App";
import DashBoardMain from "./pages/dashboard";

export const router = createBrowserRouter(createRoutesFromElements(
    <Route path="/" element={<App/>}>
        <Route path="/" element={<DashBoardMain/>}/>
        <Route path="auth/login" element={<LoginPage/>}/>
    </Route>
));
