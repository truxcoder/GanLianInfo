package controllers

import (
	"GanLianInfo/models"
	"errors"
	log "github.com/truxcoder/truxlog"
)

// 从缓存中取出身份证信息
func getIdCodeMapFromCache() map[string]string {
	_map := make(map[string]string)
	var err error
	if rdb != nil {
		exist, _ := rdb.Exists(ctx, "idCodeList").Result()
		if exist == 0 {
			setIdCodeMapToCache()
		}
		if _map, err = rdb.HGetAll(ctx, "idCodeList").Result(); err != nil {
			log.Error(err)
		}
		return _map
	}
	log.Info("redis instance is nil, now get data from database...")
	return getIdCodeMapFromDB()
}

func setIdCodeMapToCache() {
	_map := getIdCodeMapFromDB()
	if rdb != nil {
		rdb.Del(ctx, "idCodeList")
		rdb.HSet(ctx, "idCodeList", _map)
	}
}

func setDepartmentMapToCache() {
	_map := getDepartmentMapFromDB()
	if rdb != nil {
		rdb.Del(ctx, "departmentMap")
		rdb.HSet(ctx, "departmentMap", _map)
	}
}

func setDepartmentSliceToCache() {
	var err error
	var d []models.Department
	if d, err = getDepartmentSliceFromDB(); err != nil {
		return
	}
	if err = WriteDataToCache(d, "departmentSlice"); err != nil {
		log.Errorf("[setDepartmentSliceToCache error]key:departmentSlice, %v\n", err.Error())
		return
	}
}

func getDepartmentSliceFromCache() (d []models.Department, err error) {
	key := "departmentSlice"
	if err = GetDataFromCache(&d, key); err != nil {
		log.Errorf("[getDepartmentSliceFromCache]key:%s, %v\n", key, err.Error())
		if errors.Is(err, RedisKeyNotFoundERR) {
			setDepartmentSliceToCache()
			err = GetDataFromCache(&d, key)
		}
	}
	return
}

// 从缓存取出任职信息
func getPostsFromCache() (pos []posStruct, err error) {
	if err = GetDataFromCache(&pos, "posts"); err != nil {
		log.Errorf("[getPostsFromCache error]key:posts,data:[]posStruct, %v\n", err.Error())
		return
	}
	return
}

// 将任职信息写入缓存
func setPostsToCache() (err error) {
	var (
		pos []posStruct
	)
	if pos, err = getPostsFromDB(); err != nil {
		log.Errorf("[setPostsToCache error]getPostsFromDB error: %v\n", err.Error())
		return
	}
	if err = WriteDataToCache(pos, "posts"); err != nil {
		log.Errorf("[setPostsToCache error]key:posts,data:[]posStruct, %v\n", err.Error())
		return
	}
	return
}

// 将任职信息以k:v对的形式从缓存取出。传入参数isLeader，如是，则返回领导职务经历，反之亦然
func getPostMapFromCache(isLeader bool) (postsMap map[int64][]PostLevelPosition, err error) {
	var (
		key string
	)
	if isLeader {
		key = "leaderPostMap"
	} else {
		key = "nonLeaderPostMap"
	}
	if err = GetDataFromCache(&postsMap, key); err != nil {
		log.Errorf("[getPostsFromCache error]key:%s,data:[]postsMap, %v\n", key, err.Error())
	}
	return
}

// 将任职信息以k:v对的形式生成json并写入缓存，便于用人员id取出该人员的任职经历
func setPostMapToCache() (err error) {
	//var leaderPostMap = make(map[int64][]PostLevelPosition)
	//var nonLeaderPostMap = make(map[int64][]PostLevelPosition)
	leaderPostMap, nonLeaderPostMap := getPostMapFromDB()
	if err = WriteDataToCache(leaderPostMap, "leaderPostMap"); err != nil {
		log.Errorf("[setPostMapToCache error] %v\n", err.Error())
		return
	}
	if err = WriteDataToCache(nonLeaderPostMap, "nonLeaderPostMap"); err != nil {
		log.Errorf("[setPostMapToCache error] %v\n", err.Error())
		return
	}
	return
}

func CacheTest() {
	err := setPostsToCache()
	if err != nil {
		return
	}
}
