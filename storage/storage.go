package storage

type Storage interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Incr(key string) error
}
