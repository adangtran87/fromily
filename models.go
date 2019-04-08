package main

import (
	"time"
)

type RoleInfo struct {
	GuildID string
	ID      string
	Name    string
	Count   uint
	Time    time.Time
}
