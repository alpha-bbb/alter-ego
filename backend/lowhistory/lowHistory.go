package recenthistories

import "github.com/alpha-bbb/alter-ego/backend/entity"

func GetRecentHistories(histories []*entity.TalkHistory, limit int) []*entity.TalkHistory{
    if len(histories) > limit {
        return histories[len(histories)-limit:]
    }
    return histories
}
