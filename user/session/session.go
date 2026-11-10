package session

import (
	"crypto/rand"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/govwa/util/config"
	"log"
	"net/http"
)

type Self struct{}

func New() *Self {
	return &Self{}
}

var store = sessions.NewCookieStore(sessionKey())

// sessionKey returns the cookie signing key. An empty key or the published
// example key is never used; a random per-process key is generated instead.
func sessionKey() []byte {
	k := config.Cfg.Sessionkey
	if k != "" && k != "G0Vw444" {
		return []byte(k)
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func (self *Self) SetSession(w http.ResponseWriter, r *http.Request, data map[string]string) {
	session, err := store.Get(r, "govwa")

	if err != nil {
		log.Println(err.Error())
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false, //set to false for xss :)
	}

	session.Values["govwa_session"] = true

	//create new session to store on server side
	if data != nil {
		for key, value := range data {
			session.Values[key] = value
		}
	}
	err = session.Save(r, w) //safe session and send it to client as cookie

	if err != nil {
		log.Println(err.Error())
	}
}

func (self *Self) GetSession(r *http.Request, key string) string {
	session, err := store.Get(r, "govwa")

	if err != nil {
		log.Println(err.Error())
		return ""
	}
	data := session.Values[key]
	sv := fmt.Sprintf("%v", data)
	return sv
}

func (self *Self) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "govwa")
	if err != nil {
		log.Println(err.Error())
	}

	session.Options = &sessions.Options{
		MaxAge:   -1,
		HttpOnly: false, //set to false for xss :)
	}

	session.Values["govwa_session"] = false
	err = session.Save(r, w) //safe session and send it to client as cookie

	if err != nil {
		log.Println(err.Error())
	}

	return
}

func (self *Self) IsLoggedIn(r *http.Request) bool {
	s, err := store.Get(r, "govwa")
	if err != nil {
		log.Println(err.Error())
	}
	if auth, ok := s.Values["govwa_session"].(bool); !ok || !auth {
		return false
	}
	return true
}
