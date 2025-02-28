import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { Home } from "./routes/Home";
import { Privacy } from "./routes/Privacy";
import { Terms } from "./routes/Terms";
import { Tokushoho } from "./routes/Tokushoho";

const router = createBrowserRouter([
  { path: "/", element: <Home /> },
  { path: "/tokushoho", element: <Tokushoho /> },
  { path: "/privacy", element: <Privacy /> },
  { path: "/terms", element: <Terms /> },
]);

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
