package entity

type TalkHistory struct {
	Date    string // ISO 8601 date format (2021-05-22T00:00:00Z)
	User    User   // 発言したユーザー
	Message string // 発言内容
}
