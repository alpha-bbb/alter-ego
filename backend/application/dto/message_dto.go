package dto

type MessageDTO struct {
	Messages []string
	Status   int32 // TODO: 後で列挙型にする？
}
