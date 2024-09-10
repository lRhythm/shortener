package service

type RepositoryInterface interface {
	Ping() (err error)
	Put(key, value string) (err error)
	Get(key string) (value string, err error)
	Close() (err error)
}

type Client struct {
	storage RepositoryInterface
}
