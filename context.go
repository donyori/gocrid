package gocrid

import "net/http"

type WriteRespFunc func(http.ResponseWriter) bool

type Context struct {
	manager *Manager

	req *http.Request
	rw  http.ResponseWriter

	respCookie *http.Cookie

	id    string
	uInfo *userInfo

	beforeQueue []WriteRespFunc
	afterQueue  []WriteRespFunc
}

func (c *Context) GetRequest() *http.Request {
	if c == nil {
		return nil
	}
	return c.req
}

func (c *Context) IsLogin() bool {
	return c != nil && c.manager != nil && c.uInfo != nil
}

func (c *Context) GetId() (id string, err error) {
	if c.uInfo == nil {
		return "", ErrNotLogin
	}
	return c.id, nil
}

func (c *Context) GetUsername() (username string, err error) {
	if c.uInfo == nil {
		return "", ErrNotLogin
	}
	return c.uInfo.username, nil
}

func (c *Context) GetHost() (host string, err error) {
	if c.uInfo == nil {
		return "", ErrNotLogin
	}
	return c.manager.getHost(c.uInfo)
}

func (c *Context) Login(username string) error {
	if c.uInfo != nil {
		return ErrAlreadyLogin
	}
	id, info, err := c.manager.login(username, c.req)
	if err != nil {
		return err
	}
	c.id, c.uInfo = id, info
	c.respCookie = c.manager.newCookie(c.id)
	return nil
}

func (c *Context) Logout() error {
	if c.uInfo == nil {
		return ErrNotLogin
	}
	c.manager.logout(c.id, c.uInfo)
	c.respCookie = c.manager.newCookieForDelete()
	c.id = ""
	c.uInfo = nil
	return nil
}

func (c *Context) ResetTimer() error {
	if c.uInfo == nil {
		return ErrNotLogin
	}
	c.manager.resetTimer(c.id, c.uInfo)
	c.respCookie = c.manager.newCookie(c.id)
	return nil
}

// Write response before writing the Cookie.
func (c *Context) BeforeWriteResp(f WriteRespFunc) {
	c.beforeQueue = append(c.beforeQueue, f)
}

// Write response after writing the Cookie.
func (c *Context) AfterWriteResp(f WriteRespFunc) {
	c.afterQueue = append(c.afterQueue, f)
}

func (c *Context) write() {
	for _, wrf := range c.beforeQueue {
		if !wrf(c.rw) {
			return
		}
	}
	if c.respCookie != nil {
		http.SetCookie(c.rw, c.respCookie)
	}
	for _, wrf := range c.afterQueue {
		if !wrf(c.rw) {
			return
		}
	}
}
