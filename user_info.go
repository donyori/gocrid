package gocrid

import "time"

type userInfo struct {
	username     string
	timer        *time.Timer
	cleanFinishC <-chan struct{}
	host         string
}
