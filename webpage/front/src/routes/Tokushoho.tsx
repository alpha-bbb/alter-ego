import type { FC } from "react";
import "../index.css";
import { Link } from "react-router-dom";

export const Tokushoho: FC = () => {
  return (
    <div style={{ textAlign: "left" }}>
      <h1>特定商取引法に基づく表記</h1>
      <div className="mt-4">
        <p>
          ■販売業者
          <br />
          藪田 亘亨
        </p>
        <p>
          ■所在地
          <br />
          〒530-0001
          <br />
          大阪府大阪市北区梅田1丁目2番2号
          <br />
          大阪駅前第2ビル12-12
        </p>
        <p>
          ■電話番号
          <br />
          080-2875-9205
        </p>
        <p>
          ■メールアドレス
          <br />
          alterego@fuu.jp
        </p>
        <p>
          ■販売価格
          <br />
          通常プラン: 1,000円（税込）
        </p>
        <p>
          ■追加手数料等の追加料金
          <br />
          特になし
        </p>
        <p>
          ■交換および返品（返金ポリシー）
          <br />
          デジタルコンテンツの特性上、返品・キャンセルはお受けしておりません。
        </p>
        <p>
          ■引渡時期
          <br />
          決済手続き完了後、すぐに利用できるようになります。
        </p>
        <p>
          ■支払方法
          <br />
          クレジットカード決済
        </p>
        <p>
          ■決済期間
          <br />
          クレジットカード決済の場合、ご注文時にお支払いが確定します。
        </p>
        <p>
          ■サービスの提供時期
          <br />
          お支払い確認後、即時にご利用いただけます。
        </p>
        <p>
          ■動作環境
          <br />
          LINE上のBot
        </p>
      </div>
      <footer className="mt-32 text-center">
        <Link to="/" className="text-blue-500 hover:underline">
          topへ戻る
        </Link>
      </footer>
    </div>
  );
};
