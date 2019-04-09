package gocrid

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Manager struct {
	serveMux      *http.ServeMux
	settings      Settings
	idUserInfoMap sync.Map
}

func NewManager(serveMux *http.ServeMux, settings *Settings) *Manager {
	if serveMux == nil {
		serveMux = http.DefaultServeMux
	}
	if settings == nil {
		settings = NewSettings()
	}
	m := &Manager{
		serveMux: serveMux,
		settings: *settings,
	}
	m.settings.FillEmptyFields()
	return m
}

func (m *Manager) Handle(pattern string, handler Handler) {
	m.serveMux.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			// Parse request.
			ctx, err := m.parseRequest(w, r)
			// Call handler.
			handler.Handle(ctx, err)
			// Write context to response.
			ctx.write()
		})
}

func (m *Manager) HandleFunc(pattern string, handler func(*Context, error)) {
	m.Handle(pattern, HandlerFunc(handler))
}

func (m *Manager) parseRequest(w http.ResponseWriter, r *http.Request) (
	context *Context, err error) {
	context = &Context{
		manager: m,
		req:     r,
		rw:      w,
	}
	cookie, err := r.Cookie(m.settings.CookieName)
	if err != nil || cookie == nil {
		if err == http.ErrNoCookie {
			err = nil
		}
		return
	}
	id := cookie.Value
	infoItf, ok := m.idUserInfoMap.Load(id)
	if !ok {
		context.respCookie = m.newCookieForDelete()
		return
	}
	info := infoItf.(*userInfo)
	if m.settings.EnableHostBinding {
		var host string
		host, err = m.getRemoteHost(r)
		if err != nil {
			return
		}
		if host != info.host {
			context.respCookie = m.newCookieForDelete()
			return
		}
	}
	context.id = id
	context.uInfo = info
	return
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
