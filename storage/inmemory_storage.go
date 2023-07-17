package storage

import (
	"fmt"
	"strconv"
)

type inMemory struct {
	dbCount int
	storage map[int]map[string]interface{}
}

func NewInMemory(dbCntStr string) Storage {
	dbCnt, err := strconv.Atoi(dbCntStr)
	if err != nil {
		dbCnt = 16
	}

	stg := make(map[int]map[string]interface{})
	for idx := 0; idx < dbCnt; idx++ {
		stg[idx] = make(map[string]interface{})
	}

	return &inMemory{
		dbCount: dbCnt,
		storage: stg,
	}
}

func (in inMemory) Select(dbIndexStr string) (int, error) {
	dbIndex, err := strconv.Atoi(dbIndexStr)
	if err != nil {
		return 0, fmt.Errorf("(error) ERR value is not an integer or out of range")
	}
	if dbIndex < 0 || dbIndex > in.dbCount-1 {
		return 0, fmt.Errorf("(error) ERR DB index is out of range")
	}
	return dbIndex, nil
}

func (in inMemory) Set(dbIndex int, key string, value interface{}) {
	stg := in.storage[dbIndex]
	stg[key] = value
}

func (in inMemory) Get(dbIndex int, key string) interface{} {
	stg := in.storage[dbIndex]
	v, ok := stg[key]
	if !ok {
		return nil
	}
	return v
}

func (in inMemory) Del(dbIndex int, key string) interface{} {
	stg := in.storage[dbIndex]
	_, ok := stg[key]
	if !ok {
		return 0
	}
	delete(stg, key)
	return 1
}

func (in inMemory) GetAll(dbIndex int) <-chan string {
	strChan := make(chan string)
	go func() {
		for k, v := range in.storage[dbIndex] {
			strChan <- fmt.Sprintf("%s %v", k, v)
		}
		close(strChan)
	}()

	return strChan
}
