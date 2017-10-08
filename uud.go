package uud

import (
	"errors"
	"strconv"

	napping "gopkg.in/jmcvetta/napping.v3"
)

var HOST = "https://uud.com"

func VerifyAuth(code string) (LookupResponse, error) {
	res := LookupResponse{}
	r, err := napping.Post(HOST+"/verify/"+code, nil, &res, nil)
	if err == nil {
		if r.Status() > 299 {
			err = errors.New("uud returned error: " + strconv.Itoa(r.Status()))
		}
		if res.Error != "" {
			err = errors.New(res.Error)
		}
	}
	return res, err
}

func LookupUser(name string) (LookupResponse, error) {
	res := LookupResponse{}
	r, err := napping.Get(HOST+"/lookup/"+name, nil, &res, nil)
	if err == nil {
		if r.Status() > 299 {
			err = errors.New("uud returned error: " + strconv.Itoa(r.Status()))
		}
		if res.Error != "" {
			err = errors.New(res.Error)
		}
	}
	return res, err
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
