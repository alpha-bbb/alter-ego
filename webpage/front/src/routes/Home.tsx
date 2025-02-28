import type { FC } from "react";
import { Link } from "react-router-dom";
import example from "../assets/example.png";
import how from "../assets/how.png";
import solution from "../assets/solution.png";
import what from "../assets/what.png";
import "../index.css";

export const Home: FC = () => {
  return (
    <div className="space-y-12 sm:space-y-12 w-full px-4 sm:px-8 lg:px-0 flex flex-col items-center">
      <div className="w-full lg:max-w-[800px] space-y-10">
        <div>
          <img src={what} alt="What's alter-ego" className="w-full h-auto" />
        </div>
        <div className="text-center sm:text-right">
          <p className="what text-base sm:text-lg">
            <mark className="red-mark">会話に疲れた全ての現代人に捧げる</mark>
          </p>
          <p className="what mb-24 text-base sm:text-lg">
            <mark className="red-mark">究極のLINEbot</mark>
          </p>
        </div>
        <div>
          <img src={example} alt="Example use" className="w-full h-auto" />
        </div>
        <div className="text-center sm:text-right">
          <p className="example text-base sm:text-lg">
            <mark className="blue-mark">正直行きたくない誘い</mark>
          </p>
          <p className="example mb-24 text-base sm:text-lg">
            <mark className="blue-mark">
              いい感じに断って会話を終わらせたい
            </mark>
          </p>
        </div>
        <div>
          <img src={solution} alt="Solution" className="w-full h-auto" />
        </div>
        <div className="text-center">
          <p className="solution text-lg sm:text-lg">
            <mark className="yellow-mark">
              あなたの言動に合わせて返信を自動生成
            </mark>
          </p>
          <p className="text-sm sm:text-base mb-24">
            生成される返答の精度は提供いただくトーク量に依存します
          </p>
        </div>
        <div>
          <img src={how} alt="QR code" className="w-full h-auto" />
        </div>
        <div className="text-center">
          <p className="how text-base sm:text-lg">
            QRコードからalter-egoを友達登録して
          </p>
          <p className="how mb-12 text-base sm:text-lg">
            トーク履歴を送ってみましょう！
          </p>
        </div>
      </div>
      {/* <footer className="mt-48 pb-16 lg:mt-64 text-center text-base sm:text-lg flex justify-center space-x-8">
        <Link to="/tokushoho" className="text-blue-500 hover:underline">
          特定商取引法に基づく表記
        </Link>
        <Link to="/terms" className="text-blue-500 hover:underline">
          利用規約
        </Link>
        <Link to="/privacy" className="text-blue-500 hover:underline">
          プライバシーポリシー
        </Link>
      </footer> */}
      <footer className="mt-48 pb-16 lg:mt-64 text-center text-base sm:text-lg flex flex-col sm:flex-row justify-center sm:space-x-8 space-y-4 sm:space-y-0">
        <Link to="/tokushoho" className="text-blue-500 hover:underline">
          特定商取引法に基づく表記
        </Link>
        <Link to="/terms" className="text-blue-500 hover:underline">
          利用規約
        </Link>
        <Link to="/privacy" className="text-blue-500 hover:underline">
          プライバシーポリシー
        </Link>
      </footer>
    </div>
  );
};
