package gocaskdb

type store interface {
	Get(key string) string
	Set(key, value string)
	Close() bool
}
