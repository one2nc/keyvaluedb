package storage

type Storage interface {
	Select(dbIndexStr string) (int, error)
	Set(dbIndex int, key string, value interface{})
	Get(dbIndex int, key string) interface{}
	Del(dbIndex int, key string) interface{}
	GetAll(dbIndex int) <-chan string
}
