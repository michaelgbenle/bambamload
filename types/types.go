package types

import (
	"bambamload/enum"
	"time"
)

type RedisSessionInfo struct {
	Token  string            `json:"token"`
	Expiry time.Time         `json:"expiry"`
	Owner  enum.SessionOwner `json:"owner"`
	ID     string            `json:"id"`
}
