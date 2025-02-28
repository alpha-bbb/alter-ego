import type { FC } from "react";
import "../index.css";
import { Link } from "react-router-dom";

export const Privacy: FC = () => {
  return (
    <div className="w-full max-w-3xl mx-auto text-left px-4 sm:px-8">
      <h1 className="mt-8 text-center font-bold">プライバシーポリシー</h1>
      <div className="mt-4">
        <p>
          alter-ego（以下、「当サービス」といいます）は、ユーザーのプライバシーを尊重し、個人情報の適切な管理と保護に努めます。本プライバシーポリシーは、当サービスが収集する情報、利用方法、およびその管理について説明するものです。
          <br />
          <br />
          1. 収集する情報
          <br />
          当サービスは、以下の情報のみを取得・保存します。
          <br />
          ・サブスクリプションの解約に必要な情報（ユーザーIDや契約情報）
          <br />
          ・LINEのトーク履歴（LLMへの送信および会話の返信生成のため、一時的に使用し、保存は行いません）
          <br />
          <br />
          2. 情報の利用目的
          <br />
          取得した情報は、以下の目的で利用します。
          <br />
          ・サブスクリプションの解約手続きを円滑に行うため
          <br />
          ・当サービスの提供および改善のため
          <br />
          <br />
          3. 情報の管理
          <br />
          当サービスでは、取得した情報を安全に管理し、不正アクセス、紛失、破損、改ざん、漏洩などを防止するための適切な措置を講じます。
          <br />
          <br />
          4. 第三者提供について
          <br />
          当サービスは、ユーザーの個人情報を第三者に提供することはありません。
          <br />
          <br />
          5. ユーザーの権利
          <br />
          ユーザーは、当サービスに対し、以下の権利を有します。
          <br />
          ・保持されている自身の情報の確認
          <br />
          ・情報の訂正・削除の要求
          <br />
          これらの権利を行使する場合は、下記の問い合わせ先までご連絡ください。
          <br />
          <br />
          6. 問い合わせ先
          <br />
          alter-ego サポート窓口: alterego@fuu.jp
          <br />
          <br />
          7. 改定について
          <br />
          本プライバシーポリシーは、必要に応じて改定されることがあります。最新のポリシーは常に当サービス上で公開されます。
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
