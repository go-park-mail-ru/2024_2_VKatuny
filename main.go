package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/gorilla/mux"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInput struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type MyHandler struct {
	sessions map[string]uint
	users    map[string]*User
}

func NewMyHandler() *MyHandler {
	return &MyHandler{
		sessions: make(map[string]uint, 10),
		users: map[string]*User{
			"rvasily": {1, "rvasily", storage.HashPassword("love")},
		},
	}
}

func (api *MyHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	newUserInput := new(UserInput)
	decErr := decoder.Decode(newUserInput)
	fmt.Println(newUserInput)
	if decErr != nil {
		return
	}
	user, ok := api.users[newUserInput.Username]
	fmt.Println(user, ok)
	if !ok {
		http.Error(w, `no user`, 404)
		return
	}
	fmt.Println(storage.EqualHashedPasswords(storage.HashPassword("love"), "love"))
	fmt.Println(storage.EqualHashedPasswords(user.Password, newUserInput.Password), newUserInput.Password, user.Password)
	if !storage.EqualHashedPasswords(user.Password, newUserInput.Password) {
		http.Error(w, `bad pass`, 400)
		return
	}

	SID := storage.RandStringRunes(32)
	// мьютекс
	api.sessions[SID] = user.ID
	//htponly http secure
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	w.Write([]byte(SID))

}

func (api *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, `no sess`, 401)
		return
	}

	if _, ok := api.sessions[session.Value]; !ok {
		http.Error(w, `no sess`, 401)
		return
	}

	delete(api.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (api *MyHandler) Root(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	var id uint
	if err == nil && session != nil {
		id, authorized = api.sessions[session.Value]
	}

	if authorized {
		w.Write([]byte("{statusCode: 200, {id: " + strconv.Itoa(int(id)) + "}}"))
	} else {
		w.Write([]byte("{statusCode: 400}"))
	}
}

func main() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/", api.Root)
	r.HandleFunc("/login", api.Login)
	r.HandleFunc("/logout", api.Logout)

	http.ListenAndServe(":8080", r)
}
