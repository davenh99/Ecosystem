/* @refresh reload */
import { lazy } from "solid-js";
import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";

import "./index.css";
import UITesting from "./routes/UITesting";

const Home = lazy(() => import("./routes/Home"));
const Dashboard = lazy(() => import("./routes/Dashboard"));
const NotFound = lazy(() => import("./routes/NotFound"));

const root = document.getElementById("root");

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    "Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?"
  );
}

render(
  () => (
    <Router>
      <Route path="/" component={Home} />
      <Route path="/ui" component={UITesting} />
      <Route path="/dashboard" component={Dashboard} />
      <Route path="*paramName" component={NotFound} />
    </Router>
  ),
  root!
);
