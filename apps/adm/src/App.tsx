import { Outlet } from "react-router-dom";
import AuthProvider from "./lib/Auth";

function App() {
    return (
        <AuthProvider>
            <Outlet/>
        </AuthProvider>
    );
}

export default App;
