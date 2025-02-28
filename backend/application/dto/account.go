package dto

type PlatformType string

const (
	PlatformTypeUnspecified PlatformType = "PlatformTypeUnspecified"
	PlatformTypeLine        PlatformType = "PlatformTypeLine"
)

func (p PlatformType) String() string {
	return string(p)
}

type Account struct {
	PlatformType PlatformType // プラットフォームの種類
	AccountID    string       // アカウントID
}
