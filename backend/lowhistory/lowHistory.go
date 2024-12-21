package lowhistory

import "github.com/alpha-bbb/alter-ego/backend/entity"

func LowHistory(histories []*entity.TalkHistory, limit int) []*entity.TalkHistory{
    if len(histories) > limit {
        return histories[len(histories)-limit:]
    }
    return histories
}
