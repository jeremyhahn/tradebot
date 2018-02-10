package common

func (c *Context) GetUser() User {
	return c.User
}

func (c *Context) SetUser(user User) {
	c.User = user
}
