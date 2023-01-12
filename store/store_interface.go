package store

type Store interface {
	Set(key, value string) bool
	Get(key string, clear bool) string
	Verify(id, answer string, clear bool) bool
}
