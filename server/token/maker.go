package token

import "time"

type Maker interface {
	createToken(username string, duration time.Duration) (string, error)
	verifyToken(token string) (*Payload, error)
}
