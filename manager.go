package gocrid

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"
)

type userInfo struct {
	username     string
	timer        *time.Timer
	cleanFinishC <-chan struct{}
	host         string
}

type Manager struct {
	serveMux      *http.ServeMux
	settings      Settings
	idUserInfoMap sync.Map
	initLock      sync.RWMutex
}

var (
	ErrNilManager         error = errors.New("gocrid: Manager is nil")
	ErrManagerAlreadyInit error = errors.New("gocrid: Manager has already initialized")
	ErrManagerNotInit     error = errors.New("gocrid: Manager is NOT initialized")
	ErrHostBindingDisable error = errors.New("gocrid: host binding is disable")
)

func NewManager(serveMux *http.ServeMux, settings *Settings) *Manager {
	m := new(Manager)
	m.Init(serveMux, settings)
	return m
}

func (m *Manager) IsInitialized() bool {
	if m == nil {
		return false
	}
	m.initLock.RLock()
	defer m.initLock.RUnlock()
	return m.serveMux != nil
}

func (m *Manager) Init(serveMux *http.ServeMux, settings *Settings) error {
	if m == nil {
		return ErrNilManager
	}
	m.initLock.Lock()
	defer m.initLock.Unlock()
	if m.serveMux != nil {
		return ErrManagerAlreadyInit
	}
	if serveMux == nil {
		serveMux = http.DefaultServeMux
	}
	if settings == nil {
		settings = NewSettings()
	}
	m.serveMux = serveMux
	m.settings = *settings
	m.settings.FillEmptyFields()
	return nil
}

func (m *Manager) Handle(pattern string, handler Handler) error {
	if m == nil {
		return ErrNilManager
	}
	if !m.IsInitialized() {
		return ErrManagerNotInit
	}
	if handler == nil {
		return ErrNilHandler
	}
	m.serveMux.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			// Parse request.
			ctx, _ := m.ParseRequest(w, r) // Ignore error because it must be nil.
			// Call handler.
			handler.Handle(r, ctx)
			// Write context to response.
			if !ctx.IsWritten() {
				ctx.Write() // Ignore error.
			}
		})
	return nil
}

func (m *Manager) HandleFunc(pattern string,
	handler func(*http.Request, *Context)) error {
	return m.Handle(pattern, HandlerFunc(handler))
}

func (m *Manager) ParseRequest(w http.ResponseWriter, r *http.Request) (
	context *Context, err error) {
	if m == nil {
		return nil, ErrNilManager
	}
	if !m.IsInitialized() {
		return nil, ErrManagerNotInit
	}
	context = &Context{
		manager:   m,
		req:       r,
		rw:        w,
		isWritten: false,
	}
	cookie, err := r.Cookie(m.settings.CookieName)
	if err != nil || cookie == nil {
		return context, nil
	}
	err = nil
	context.reqCookie = cookie
	id := cookie.Value
	infoItf, ok := m.idUserInfoMap.Load(id)
	if !ok {
		context.respCookie = m.newCookieForDelete()
		return context, nil
	}
	info := infoItf.(*userInfo)
	if m.settings.EnableHostBinding {
		host, err := m.getRemoteHost(r)
		if err != nil {
			return nil, err
		}
		if host != info.host {
			context.respCookie = m.newCookieForDelete()
			return context, nil
		}
	}
	context.id = id
	context.uInfo = info
	return context, nil
}

func (m *Manager) getRemoteHost(r *http.Request) (host string, err error) {
	if !m.settings.EnableHostBinding {
		return "", ErrHostBindingDisable
	}
	host, _, err = net.SplitHostPort(r.RemoteAddr)
	return host, err
}

func (m *Manager) newCookie(value string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     m.settings.CookieName,
		Value:    value,
		Path:     m.settings.CookiePath,
		Domain:   m.settings.CookieDomain,
		Secure:   m.settings.CookieSecure,
		HttpOnly: m.settings.CookieHttpOnly,
	}
	maxAge := m.settings.MaxAge.Round(time.Second)
	if maxAge > 0 {
		cookie.Expires = time.Now().Add(maxAge)
		cookie.MaxAge = int(maxAge.Seconds())
	}
	return cookie
}

func (m *Manager) newCookieForDelete() *http.Cookie {
	cookie := &http.Cookie{
		Name:    m.settings.CookieName,
		Value:   "",
		Path:    m.settings.CookiePath,
		Domain:  m.settings.CookieDomain,
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	return cookie
}

func (m *Manager) getHost(info *userInfo) (host string, err error) {
	if !m.settings.EnableHostBinding {
		return "", ErrHostBindingDisable
	}
	return info.host, nil
}

func (m *Manager) newTimer(id string) (timer *time.Timer,
	cleanFinishC <-chan struct{}) {
	if m.settings.MaxAge <= 0 {
		return nil, nil
	}
	cleanFinishChan := make(chan struct{}, 1)
	timer = time.AfterFunc(m.settings.MaxAge, func() {
		m.idUserInfoMap.Delete(id)
		cleanFinishChan <- struct{}{}
	})
	return timer, cleanFinishChan
}

func (m *Manager) resetTimer(id string, info *userInfo) {
	if m.settings.MaxAge > 0 {
		if info.timer != nil {
			if !info.timer.Stop() {
				<-info.cleanFinishC
				m.idUserInfoMap.Store(id, info)
			}
			info.timer.Reset(m.settings.MaxAge)
		} else {
			info.timer, info.cleanFinishC = m.newTimer(id)
		}
	} else {
		if info.timer != nil {
			if !info.timer.Stop() {
				<-info.cleanFinishC
				m.idUserInfoMap.Store(id, info)
			}
			info.timer = nil
			info.cleanFinishC = nil
		}
	}
}

func (m *Manager) login(username string, r *http.Request) (
	id string, info *userInfo, err error) {
	info = &userInfo{username: username}
	if m.settings.EnableHostBinding {
		info.host, err = m.getRemoteHost(r)
		if err != nil {
			return "", nil, err
		}
	}
	prefix := m.settings.IdPrefix
	if prefix != "" {
		prefix += "-"
	}
	var newBody string
	next := true
	for next {
		newBody, err = randString(m.settings.IdBodyLength)
		if err != nil {
			return "", nil, err
		}
		id = prefix + newBody
		_, next = m.idUserInfoMap.LoadOrStore(id, info)
	}
	if m.settings.MaxAge > 0 {
		info.timer, info.cleanFinishC = m.newTimer(id)
	}
	return id, info, nil
}

func (m *Manager) logout(id string, info *userInfo) {
	if info.timer != nil && !info.timer.Stop() {
		<-info.cleanFinishC
	} else {
		m.idUserInfoMap.Delete(id)
	}
}
