package domain

type Cache interface {
	Get(key string) (CacheEntry, bool)
	Set(key string, val CacheEntry) error
	Delete(key string) error
}

type EventLoop interface {
	Run()
	Send(event Event)
}
