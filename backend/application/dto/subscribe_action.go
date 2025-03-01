package dto

// SubscribeAction はサブスクアクションの種類を表す列挙型です。
// proto の SubscribeAction enum と対応します。
type SubscribeAction int32

const (
	SubscribeActionUnspecified SubscribeAction = 0 // 未定義
	SubscribeActionStart       SubscribeAction = 1 // サブスク開始
	SubscribeActionCancel      SubscribeAction = 2 // サブスクキャンセル
	SubscribeActionCheck       SubscribeAction = 3 // サブスクチェック
)

// String は SubscribeAction の文字列表現を返します。
func (s SubscribeAction) String() string {
	switch s {
	case SubscribeActionStart:
		return "SUBSCRIBE_ACTION_START"
	case SubscribeActionCancel:
		return "SUBSCRIBE_ACTION_CANCEL"
	case SubscribeActionCheck:
		return "SUBSCRIBE_ACTION_CHECK"
	default:
		return "SUBSCRIBE_ACTION_UNSPECIFIED"
	}
}
