package auth

import "net/http"

// Strategy is the auth interface.
type Strategy interface {
	Apply(req *http.Request)
}

type NoAuth struct{}

func (NoAuth) Apply(req *http.Request) {}

type Basic struct {
	User string
	Pass string
}

func (b Basic) Apply(req *http.Request) {
	req.SetBasicAuth(b.User, b.Pass)
}

type Bearer struct {
	Token string
}

func (b Bearer) Apply(req *http.Request) {
	if b.Token != "" {
		req.Header.Set("Authorization", "Bearer "+b.Token)
	}
}
