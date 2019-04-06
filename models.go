package main

import (
	"time"
)

type Configs struct {
	Token       string   `json:"TOKEN"`
	Prefix      string   `json:"PREFIX"`
	AdminPrefix string   `json:"ADMIN_PREFIX"`
	Admins      []string `json:"ADMINS"`
}

type RoleInfo struct {
	GuildID string
	ID      string
	Name    string
	Count   uint
	Time    time.Time
}
