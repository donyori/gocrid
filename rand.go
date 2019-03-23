package gocrid

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

const charset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789"

var (
	ErrRandExit error = errors.New("gocrid: rand is exit")

	randByteC <-chan byte
	randExitC chan<- struct{}
	randLock  sync.RWMutex
)

func init() {
	RandStart()
}

func RandStart() {
	randLock.Lock()
	defer randLock.Unlock()
	if randExitC != nil {
		return
	}
	rbc := make(chan byte, 64)
	rec := make(chan struct{})
	go func(outC chan<- byte, exitC <-chan struct{}) {
		defer close(outC)
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		length := int32(len(charset))
		isExit := false
		for !isExit {
			select {
			case <-exitC:
				isExit = true
			default:
				idx := random.Int31n(length)
				outC <- charset[idx]
			}
		}
	}(rbc, rec)
	randByteC = rbc
	randExitC = rec
}

func RandExit() {
	randLock.Lock()
	defer randLock.Unlock()
	if randExitC == nil {
		return
	}
	close(randExitC)
	randExitC = nil
	for _ = range randByteC {
		// Drain randByteC.
	}
}

func randString(length int) (s string, err error) {
	if length <= 0 {
		return "", nil
	}
	randLock.RLock()
	defer randLock.RUnlock()
	b := make([]byte, length)
	var ok bool
	for i := 0; i < length; i++ {
		b[i], ok = <-randByteC
		if !ok {
			return "", ErrRandExit
		}
	}
	return string(b), nil
}
