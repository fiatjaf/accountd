package accountd

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"

	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

var HOST = "https://accountd.xyz"

func NewClient() Client {
	resp, err := http.Get(HOST + "/public-key")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	keyb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	pubkey, err := jwt.ParseRSAPublicKeyFromPEM(keyb)
	if err != nil {
		panic(err)
	}

	return Client{
		PubKey: pubkey,
	}
}

type Client struct {
	PubKey *rsa.PublicKey
}

func (c Client) VerifyAuth(tokenstr string) (data TokenData, err error) {
	token, err := jwt.ParseWithClaims(
		tokenstr, &data, func(*jwt.Token) (interface{}, error) {
			return c.PubKey, nil
		},
	)
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*TokenData); ok && token.Valid {
		data = *claims
		return
	} else {
		return data, errors.New("jwt token invalid fields.")
	}
}

type TokenData struct {
	User struct {
		Name string
		Id   int
	}
	Accounts []string
	jwt.StandardClaims
}
