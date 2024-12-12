package models

import (
	"github.com/gorilla/sessions"
)

// Store is a global variable holding the session store
var Store = sessions.NewCookieStore([]byte("super-secret-key"))
