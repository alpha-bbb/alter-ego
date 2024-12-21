package functions

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func ReplacePhoneNumbers(input string) string {
    phonePattern := `\b\d{2,4}-\d{2,4}-\d{4}\b`

    re := regexp.MustCompile(phonePattern)
    input = re.ReplaceAllStringFunc(input, func(phone string) string {
        // トークン化し、メモリに保存
        token := fmt.Sprintf("［電話番号#%02d］", GetNextTokenIndex())
        tokenMap[token] = phone
        return token
    })

    return input
}

func ReplacePrefectures(input string) string {
    jsonFilePath := "./json/prefectures.json"
    var prefectures Prefectures
    if err := LoadJSON(jsonFilePath, &prefectures); err != nil {
        return input
    }

    sort.Slice(prefectures.Prefectures, func(i, j int) bool {
        return len(prefectures.Prefectures[i]) > len(prefectures.Prefectures[j])
    })

    for _, pref := range prefectures.Prefectures {
        pattern := fmt.Sprintf("(%s)(都|道|府|県)", regexp.QuoteMeta(pref))
        re := regexp.MustCompile(pattern)
        
        input = re.ReplaceAllStringFunc(input, func(matched string) string {
            token := fmt.Sprintf("［都道府県#%02d］", GetNextTokenIndex())
            tokenMap[token] = matched
            return token
        })
    }

    return input
}

type MunicipalityData struct {
	Municipalities []string `json:"municipalities"`
}

// 市区町村の置き換え処理
func ReplaceMunicipalities(input string) string {
    jsonFilePath := "./json/municipalities.json"
    var municipalityData MunicipalityData

    // JSONを読み込む
    if err := LoadJSON(jsonFilePath, &municipalityData); err != nil {
        return input
    }

    // 市区町村名で置換
    for _, muni := range municipalityData.Municipalities {
        // 市区町村名がinputに含まれているかをチェック
        if strings.Contains(input, muni) {
            // 置換処理
            token := fmt.Sprintf("［市区町村#%02d］", GetNextTokenIndex())
            input = strings.ReplaceAll(input, muni, token)
            tokenMap[token] = muni
        }
    }

    // 置換後の最終的な文字列を返す
    return input
}

// トークンを元の値に戻す処理
func RestoreTokens(input string) string {
    tokenMapMutex.Lock()
    defer tokenMapMutex.Unlock()

    // トークンを元の値に戻す
    for token, original := range tokenMap {
        input = strings.ReplaceAll(input, token, original)
    }
    return input
}


type Prefectures struct {
    Prefectures []string `json:"prefectures"`
}
