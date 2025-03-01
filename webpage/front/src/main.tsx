import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { Home } from "./routes/Home";
import { Subscribed } from "./routes/Subscribed";
import { Tokushoho } from "./routes/Tokushoho";

const router = createBrowserRouter([
  { path: "/", element: <Home /> },
  { path: "/tokushoho", element: <Tokushoho /> },
  { path: "/subscribed", element: <Subscribed /> },
]);

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
