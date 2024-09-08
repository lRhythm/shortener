package service

type repositoryInterface interface {
	Put(key, value string) (err error)
	Get(key string) (value string, err error)
	Close() (err error)
}

type dbInterface interface {
	Ping() (err error)
}

type Client struct {
	storage repositoryInterface
	db      dbInterface
}
