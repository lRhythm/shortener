/*
Package service - слой бизнес логики сервиса.
*/
package service

// New - конструктор Client.
func New(opts ...func(*Client)) *Client {
	c := new(Client)
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithStorage - добавление имплементации интерфейса RepositoryInterface в Client.
// Является аргументом функции New.
func WithStorage(i RepositoryInterface) func(*Client) {
	return func(c *Client) {
		c.storage = i
	}
}
