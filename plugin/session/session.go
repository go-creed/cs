package session

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var (
	sessionPrefix = "session-x-"
	store         *sessions.CookieStore
)

func Init() {
	store = sessions.NewCookieStore([]byte("x984r739#$%^$5egodfgosDSFASDG25@#"))
}

// Use gin's context to get session
func GetSessionGin(ctx *gin.Context) *sessions.Session {
	return GetSession(ctx.Writer, ctx.Request)
}

// Use http to get session
func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {

	var xId string
	for _, cookie := range r.Cookies() {
		if strings.Index(cookie.Name, sessionPrefix) == 0 {
			xId = cookie.Name
			break
		}
	}
	if xId == "" {
		xId = sessionPrefix + uuid.New().String()
	}
	ses, _ := store.Get(r, xId)

	if ses.ID == "" {
		// Store session id in cookie
		cookie := &http.Cookie{Name: xId, Value: xId, Path: "/", Expires: time.Now().Add(30 * time.Second), MaxAge: 0}
		http.SetCookie(w, cookie)
		// Store new session
		ses.ID = xId
		ses.Save(r, w)
	}

	return ses
}
