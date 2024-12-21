package preprocess

import (
	"github.com/alpha-bbb/alter-ego/backend/preprocess/functions"
)

// メッセージを一度にトークン化する関数
func TokenizeMessage(message string) string {
    message = functions.ReplacePhoneNumbers(message)
    message = functions.ReplacePrefectures(message)
    message = functions.ReplaceMunicipalities(message)
    return message
}

func RestoreTokens(input string) string {
    input = functions.RestoreTokens(input)
    return input
}

func ClearAllTokens() {
    functions.ClearAllTokens()
}
