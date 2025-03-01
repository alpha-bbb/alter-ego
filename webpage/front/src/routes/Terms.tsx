import type { FC } from "react";
import "../index.css";
import { Link } from "react-router-dom";

export const Terms: FC = () => {
  return (
    <div className="w-full max-w-3xl mx-auto text-left px-4 sm:px-8">
      <h1 className="mt-8 text-center font-bold">利用規約</h1>
      <div className="mt-4">
        <p>
          alter-ego（以下、「当サービス」といいます）を利用する際の条件を定めるものです。
          <br />
          <br />
          1. 適用範囲
          <br />
          本規約は、日本国内のすべてのユーザーに適用されます。
          <br />
          <br />
          2. サービスの内容
          <br />
          当サービスは、ユーザーのLINEのトーク履歴を元に、LLM（大規模言語モデル）を利用して会話の返信を生成するLINE
          botです。
          <br />
          <br />
          3. 禁止行為
          <br />
          ユーザーは、以下の行為を行わないものとします。
          <br />
          ・違法行為または公序良俗に反する行為
          <br />
          ・当サービスの運営を妨害する行為
          <br />
          ・その他、当サービスが不適切と判断する行為
          <br />
          <br />
          4. 免責事項
          <br />
          当サービスは、生成された返信の内容の正確性や適切性について保証しません。
          <br />
          当サービスの利用により発生したいかなる損害についても、当サービスは一切の責任を負いません。
          <br />
          <br />
          5. サービスの変更・停止
          <br />
          当サービスは、ユーザーへの事前通知なしに、サービスの内容を変更・停止することがあります。
          <br />
          <br />
          6. 利用規約の変更
          <br />
          本規約は、必要に応じて改定されることがあります。改定後の規約は、当サービス上で公開された時点で効力を持つものとします。
          <br />
          <br />
          7. 問い合わせ先
          <br />
          alter-ego サポート窓口: alterego@fuu.jp
        </p>
      </div>
      <div className="mt-8 text-right text-sm text-gray-500">
        公開日: 2025年2月28日
      </div>
      <footer className="mt-16 mb-16 text-center">
        <Link to="/" className="text-blue-500 hover:underline">
          TOPへ戻る
        </Link>
      </footer>
    </div>
  );
};
