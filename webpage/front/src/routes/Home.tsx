import type { FC } from "react";
import { Link } from "react-router-dom";
import example from "../assets/example.jpg";
import solution from "../assets/solution.jpg";
import what from "../assets/what.jpg";
import "../index.css";

export const Home: FC = () => {
  return (
    <div
      className="space-y-6 gap-12">
      <div className="flex flex-col items-center space-y-10 w-full max-w-[600px]">
        <div><img src={what} alt="What's alter-ego" className="w-full h-auto" /></div>
        <div><img src={example} alt="Example usecase" className="w-full h-auto" /></div>
        <div><img src={solution} alt="Solution" className="w-full h-auto" /></div>
      </div>
      <footer className="mt-32 text-center">
        <Link to="/tokushoho" className="text-blue-500 hover:underline">
          特定商取引法に基づく表記
        </Link>
      </footer>
    </div>
  );
};
