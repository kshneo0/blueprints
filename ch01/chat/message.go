package main

import "time"

//message는 단일 메시지를 나타낸다
type message struct {
	Name string
	Message string
	When time.Time
	AvatarURL string
}