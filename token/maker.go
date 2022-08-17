package token

import "time"

//Maker is an interface for managing tokens
type Maker interface {
	//CreateToken create and sign a new token for a specific user
	CreateToken(username string, duration time.Duration) (string, error)

	//VerifyToken verify if the token is valid
	VerifyToken(token string) (*Payload, error)
}
