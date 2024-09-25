package murl

import (
	"fmt"
	"net/url"
)

type Mutation func(*url.URL)

func Mutate(rawURL string, actions ...Mutation) string {
	u, _ := url.Parse(rawURL)
	for _, action := range actions {
		action(u)
	}
	return u.String()
}

func AddQuery(keysAndValues ...any) Mutation {
	return func(u *url.URL) {
		q := u.Query()
		for i := 0; i < len(keysAndValues); i += 2 {
			q.Set(keysAndValues[i].(string), fmt.Sprint(keysAndValues[i+1]))
		}
		u.RawQuery = q.Encode()
	}
}

func RemoveQuery(keys ...string) Mutation {
	return func(u *url.URL) {
		q := u.Query()
		for _, key := range keys {
			q.Del(key)
		}
		u.RawQuery = q.Encode()
	}
}
