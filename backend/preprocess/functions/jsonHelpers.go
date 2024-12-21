package functions

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONファイルを読み込む
func LoadJSON(filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ファイルを開けません: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("JSON デコードエラー: %w", err)
	}

	return nil
}
