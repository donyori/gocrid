package gocrid

import (
	"errors"
	"net/http"
	"sync"
)

var (
	ErrNilContext            error = errors.New("gocrid: Context is nil")
	ErrInvalidContext        error = errors.New("gocrid: Context is invalid")
	ErrContextAlreadyWritten error = errors.New("gocrid: Context has already written")

	ErrAlreadyLogin error = errors.New("gocrid: user has already login")
	ErrNotLogin     error = errors.New("gocrid: user is NOT login")
)

type WriteRespFunc func(http.ResponseWriter) bool

type Context struct {
	manager *Manager

	req *http.Request
	rw  http.ResponseWriter

	reqCookie  *http.Cookie
	respCookie *http.Cookie

	id    string
	uInfo *userInfo

	isWritten            bool
	beforeWriteRespQueue []WriteRespFunc
	afterWriteRespQueue  []WriteRespFunc
	writeLock            sync.RWMutex
}

func (c *Context) IsLogin() bool {
	return c != nil && c.manager != nil && c.uInfo != nil
}

func (c *Context) GetId() (id string, err error) {
	if c == nil {
		return "", ErrNilContext
	}
	if c.manager == nil {
		return "", ErrInvalidContext
	}
	if c.uInfo == nil {
		return "", ErrNotLogin
	}
	return c.id, nil
}

func (c *Context) GetUsername() (username string, err error) {
	if c == nil {
		return "", ErrNilContext
	}
	if c.manager == nil {
		return "", ErrInvalidContext
	}
	if c.uInfo == nil {
		return "", ErrNotLogin
	}
	return c.uInfo.username, nil
}

func (c *Context) GetHost() (host string, err error) {
	if c == nil {
		return "", ErrNilContext
	}
	if c.manager == nil {
		return "", ErrInvalidContext
	}
	if c.uInfo == nil {
		return "", ErrNotLogin
	}
	return c.manager.getHost(c.uInfo)
}

func (c *Context) Login(username string) error {
	if c == nil {
		return ErrNilContext
	}
	if c.manager == nil {
		return ErrInvalidContext
	}
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
	if c == nil {
		return ErrNilContext
	}
	if c.manager == nil {
		return ErrInvalidContext
	}
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
	if c == nil {
		return ErrNilContext
	}
	if c.manager == nil {
		return ErrInvalidContext
	}
	if c.uInfo == nil {
		return ErrNotLogin
	}
	c.manager.resetTimer(c.id, c.uInfo)
	c.respCookie = c.manager.newCookie(c.id)
	return nil
}

func (c *Context) IsWritten() bool {
	if c == nil || c.manager == nil {
		return false
	}
	c.writeLock.RLock()
	defer c.writeLock.RUnlock()
	return c.isWritten
}

func (c *Context) BeforeWriteResp(f func(http.ResponseWriter) bool) error {
	if c == nil {
		return ErrNilContext
	}
	if c.manager == nil {
		return ErrInvalidContext
	}
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	if c.isWritten {
		return ErrContextAlreadyWritten
	}
	c.beforeWriteRespQueue = append(c.beforeWriteRespQueue, WriteRespFunc(f))
	return nil
}

func (c *Context) AfterWriteResp(f func(http.ResponseWriter) bool) error {
	if c == nil {
		return ErrNilContext
	}
	if c.manager == nil {
		return ErrInvalidContext
	}
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	if c.isWritten {
		return ErrContextAlreadyWritten
	}
	c.afterWriteRespQueue = append(c.afterWriteRespQueue, WriteRespFunc(f))
	return nil
}

func (c *Context) Write() error {
	if c == nil {
		return ErrNilContext
	}
	if c.manager == nil {
		return ErrInvalidContext
	}
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	if c.isWritten {
		return ErrContextAlreadyWritten
	}
	defer func() {
		c.isWritten = true
	}()
	for _, wrf := range c.beforeWriteRespQueue {
		if !wrf(c.rw) {
			return nil
		}
	}
	if c.respCookie != nil {
		http.SetCookie(c.rw, c.respCookie)
	}
	for _, wrf := range c.afterWriteRespQueue {
		if !wrf(c.rw) {
			return nil
		}
	}
	return nil
}
