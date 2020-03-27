package sessions

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func GetSession(r *http.Request,name_cookie string)(name string){
	if cookie,err:=r.Cookie(name_cookie);err==nil{
		cookieValue:= make(map[string]string)
		if err=cookieHandler.Decode(name_cookie,cookie.Value,&cookieValue);err==nil{
			name = cookieValue[name_cookie]
		}
	}
	return name
}
func SetSession(key_cookie,value_cookie string,w http.ResponseWriter){
	value:=map[string]string{
		key_cookie:value_cookie,
	}
	if encoded,err:=cookieHandler.Encode(key_cookie,value);err==nil{
		cookie:=&http.Cookie{
			Name:       key_cookie,
			Value:      encoded,
			Path:       "/",
			MaxAge:     3600,
		}
		http.SetCookie(w,cookie)
	}
}

func ClearSession(name string,w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
