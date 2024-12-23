package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/big-dust/DreamBridge/internal/model/major"
    "github.com/go-redis/redis/v8"
)

const (
    // 缓存键前缀
    schoolScoresPrefix = "school:scores:"
    // 缓存过期时间
    schoolScoresTTL = 24 * time.Hour
)

// GetSchoolScoresCache 从缓存获取学校分数信息
func GetSchoolScoresCache(schoolName string, year int, studentType string) ([]major.SchoolScore, error) {
    key := fmt.Sprintf("%s%s:%d:%s", schoolScoresPrefix, schoolName, year, studentType)
    
    // 尝试从缓存获取
    val, err := RedisClient.Get(context.Background(), key).Result()
    if err == redis.Nil {
        // 缓存未命中，从数据库查询
        scores, err := major.GetSchoolScores(schoolName, year, studentType)
        if err != nil {
            return nil, err
        }

        // 存入缓存
        data, err := json.Marshal(scores)
        if err != nil {
            return nil, err
        }
        
        err = RedisClient.Set(context.Background(), key, data, schoolScoresTTL).Err()
        if err != nil {
            // 缓存存储失败只记录日志，不影响返回
            log.Printf("存储缓存失败: %v", err)
        }

        return scores, nil
    } else if err != nil {
        return nil, err
    }

    // 缓存命中，解析数据
    var scores []major.SchoolScore
    err = json.Unmarshal([]byte(val), &scores)
    if err != nil {
        return nil, err
    }

    return scores, nil
}

// ClearSchoolScoresCache 清除学校分数缓存
func ClearSchoolScoresCache(schoolName string, year int, studentType string) error {
    key := fmt.Sprintf("%s%s:%d:%s", schoolScoresPrefix, schoolName, year, studentType)
    return RedisClient.Del(context.Background(), key).Err()
} 