import type { FC } from "react";
import "../index.css";
import { Link } from "react-router-dom";

export const Tokushoho: FC = () => {
  return (
    <div>
      <h1>特定商取引法に基づく表記</h1>
      <p>内容</p>
      <footer className="mt-32 text-center">
        <Link to="/" className="text-blue-500 hover:underline">
          topへ戻る
        </Link>
      </footer>
    </div>
  );
};
