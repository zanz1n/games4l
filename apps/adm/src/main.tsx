import { Route, Router, useRouter } from "preact-router";
import "./index.css";
import { render } from "preact";
import { useEffect } from "preact/hooks";
import LoginPage from "./pages/LoginPage";
import DashBoardMain from "./pages/dashboard";

function RedirectTo(to: string) {
    return () => {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const [_, navigate] = useRouter();

        useEffect(() => {
            navigate(to);
        }, []);

        return <h1 style={{
            textAlign: "center",
            marginTop: "32px"
        }}>Redirecting...</h1>;
    };
}

function Main() {
    return <Router>
        <Route path="/" component={RedirectTo("/dash")} />
        <Route path="/auth/login" component={LoginPage} />
        <Route path="/dash" component={DashBoardMain} />
        <Route
            path="/dash/create-question"
            component={DashBoardMain}
        />
        <Route
            path="/dash/update-question"
            component={DashBoardMain}
        />
    </Router>;
}

render(<Main />, document.getElementById("root") as HTMLElement);
