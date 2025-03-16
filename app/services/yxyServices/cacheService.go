package yxyServices

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	r "wejh-go/config/redis"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
)

func GetElecRoomStrConcat(yxyUid, campus string) (*string, error) {
	cacheKey := "elec:room_str_concat:" + campus + ":" + yxyUid
	cachedRoomStrConcat, err := r.RedisClient.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		balance, err := ElectricityBalance(yxyUid, campus)
		if err != nil {
			return nil, err
		}
		err = r.RedisClient.Set(ctx, cacheKey, balance.RoomStrConcat, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}
		return &balance.RoomStrConcat, nil
	} else if err != nil {
		return nil, err
	}
	return &cachedRoomStrConcat, nil
}

func GetElecConsumptionRecords(yxyUid, campus, roomStrConcat string) (*EleConsumptionRecords, error) {
	cacheKey := "elec:consumption_records:" + roomStrConcat
	cachedRecords, err := r.RedisClient.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		records, err := ElectricityConsumptionRecords(yxyUid, campus, roomStrConcat)
		if err != nil {
			return nil, err
		}
		recordsJSON, err := json.Marshal(records)
		if err != nil {
			return nil, err
		}
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		ttl := time.Until(midnight)
		err = r.RedisClient.Set(ctx, cacheKey, recordsJSON, ttl).Err()
		if err != nil {
			return nil, err
		}
		return records, nil
	} else if err != nil {
		return nil, err
	}
	var records EleConsumptionRecords
	err = json.Unmarshal([]byte(cachedRecords), &records)
	if err != nil {
		return nil, err
	}
	return &records, nil
}

func GetBusAuthToken(yxyUid string) (*string, error) {
	cacheKey := "bus:auth_token:" + yxyUid
	cachedToken, err := r.RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// 使用 singleflight 防止缓存击穿
		token, err, _ := g.Do(cacheKey, func() (interface{}, error) {
			t, e := BusAuth(yxyUid)
			if e != nil {
				return nil, e
			}
			e = r.RedisClient.Set(ctx, cacheKey, *t, 7*24*time.Hour).Err()
			if e != nil {
				return nil, e
			}
			return t, nil
		})
		if err != nil {
			return nil, err
		}
		return token.(*string), nil
	} else if err != nil {
		return nil, err
	}
	return &cachedToken, nil
}

// func GetBusInfo(page, pageSize, search string) (*busResp, error) {
// 	cacheKey := "bus:info:"
// 	cachedInfo, err := r.RedisClient.Get(ctx, cacheKey).Result()
// 	if err == redis.Nil {
// 		info, err := BusInfo(page, pageSize, search)
// 		if err != nil {
// 			return nil, err
// 		}
// 		infoJSON, err := json.Marshal(info)
// 		if err != nil {
// 			return nil, err
// 		}
// 		now := time.Now()
// 		midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
// 		ttl := time.Until(midnight)
// 		err = r.RedisClient.Set(ctx, cacheKey, infoJSON, ttl).Err()
// 		if err != nil {
// 			return nil, err
// 		}
// 		return info, nil
// 	} else if err != nil {
// 		return nil, err
// 	}
// 	var info busResp
// 	err = json.Unmarshal([]byte(cachedInfo), &info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &info, nil
// }
