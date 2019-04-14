package main

import (
	"fmt"
)

var Responses = map[string]string{
	"cmd:unknown":     ":upside_down: Unknown Command: `%s`",
	"cmd:unknown_sub": ":upside_down: Unknown Sub Command: `%s`",
}

//TODO: This assumes that all Responses need an argument. Need to figure out a
//better way to handle this
func GetResp(key string, cmdAttempt string) string {
	return fmt.Sprintf(Responses[key], cmdAttempt)
}
