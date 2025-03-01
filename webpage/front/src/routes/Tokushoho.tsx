import type { FC } from "react";
import "../index.css";
import { Link } from "react-router-dom";

export const Tokushoho: FC = () => {
  return (
    <div className="w-full max-w-3xl mx-auto text-left px-4 sm:px-8">
      <h1 className="mt-8 text-center font-bold">特定商取引法に基づく表記</h1>
      <div className="mt-4 space-y-4">
        <p>
          ■販売業者
          <br />
          <span className="ml-4">藪田 亘亨</span>
        </p>
        <p>
          ■所在地
          <br />
          <span className="ml-4">〒530-0001</span>
          <br />
          <span className="ml-4">大阪府大阪市北区梅田1丁目2番2号</span>
          <br />
          <span className="ml-4">大阪駅前第2ビル12-12</span>
        </p>
        <p>
          ■電話番号
          <br />
          <span className="ml-4">080-2875-9205</span>
        </p>
        <p>
          ■メールアドレス
          <br />
          <span className="ml-4">alterego@fuu.jp</span>
        </p>
        <p>
          ■販売価格
          <br />
          <span className="ml-4">通常プラン: 1,000円（税込）</span>
        </p>
        <p>
          ■追加手数料等の追加料金
          <br />
          <span className="ml-4">特になし</span>
        </p>
        <p>
          ■交換および返品（返金ポリシー）
          <br />
          <span className="ml-4">
            デジタルコンテンツの特性上、返品・キャンセルはお受けしておりません。
          </span>
        </p>
        <p>
          ■引渡時期
          <br />
          <span className="ml-4">
            決済手続き完了後、すぐに利用できるようになります。
          </span>
        </p>
        <p>
          ■支払方法
          <br />
          <span className="ml-4">クレジットカード決済</span>
        </p>
        <p>
          ■決済期間
          <br />
          <span className="ml-4">
            クレジットカード決済の場合、ご注文時にお支払いが確定します。
          </span>
        </p>
        <p>
          ■サービスの提供時期
          <br />
          <span className="ml-4">
            お支払い確認後、即時にご利用いただけます。
          </span>
        </p>
        <p>
          ■動作環境
          <br />
          <span className="ml-4">LINE上のBot</span>
        </p>
      </div>
      <footer className="mt-16 mb-16 text-center">
        <Link to="/" className="text-blue-500 hover:underline">
          TOPへ戻る
        </Link>
      </footer>
    </div>
  );
};
