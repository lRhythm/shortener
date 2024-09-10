package service

func New(opts ...func(*Client)) *Client {
	c := new(Client)
	for _, o := range opts {
		o(c)
	}
	return c
}

func WithStorage(i RepositoryInterface) func(*Client) {
	return func(c *Client) {
		c.storage = i
	}
}
