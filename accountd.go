package accountd

import (
	"crypto/rsa"
	"encoding/json"
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

func (c Client) VerifyAuth(tokenstr string) (user string, err error) {
	token, err := jwt.Parse(tokenstr, func(*jwt.Token) (interface{}, error) {
		return c.PubKey, nil
	})
	if err != nil {
		return "", err
	}
	iuser, ok := token.Claims.(jwt.MapClaims)["user"]
	if !ok {
		return "", errors.New("jwt token doesn't have field 'user'.")
	}

	user, ok = iuser.(string)
	if !ok || user == "" {
		return "", errors.New("'user' field doesn't contain a string.")
	}
	return
}

func Lookup(name string) (res LookupResponse, err error) {
	resp, err := http.Get(HOST + "/lookup/" + name)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	return
}

type LookupResponse struct {
	// this will be returned only if the account
	// used in the lookup is not known
	Type string `json:"type"`

	Id       string `json:"id"`
	Accounts []struct {
		Type    string `json:"type"`
		Account string `json:"account"`
	} `json:"accounts"`

	Error string `json:"error"`
}
