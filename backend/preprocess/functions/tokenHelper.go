package functions

import "sync"

// グローバルマップとミューテックスを用意
var (
	tokenMap      = make(map[string]string) // トークンと元の情報を保存
	tokenIndex    = 0                       // トークンのインデックス
	tokenMapMutex sync.Mutex                // 並行処理用のロック
)

// 次のトークンのインデックスを取得
func GetNextTokenIndex() int {
	tokenMapMutex.Lock()
	defer tokenMapMutex.Unlock()

	index := tokenIndex
	tokenIndex++
	return index
}

// すべてのトークンを削除する関数
func ClearAllTokens() {
    tokenMapMutex.Lock()
    defer tokenMapMutex.Unlock()

    tokenMap = make(map[string]string)
}
