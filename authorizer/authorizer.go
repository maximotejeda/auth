// here go the handlers for the service

package authorizer

import "net/http"

// Login
// Method POST
func Login(w http.ResponseWriter, r *http.Request) {}

// Register
// Method POST
func Register(w http.ResponseWriter, r *http.Request) {}

// ChangePwd
// Method GET
func ChangePwd(w http.ResponseWriter, r *http.Request) {}

// Validate
// Method Post
func Validate(w http.ResponseWriter, r *http.Request) {}

// Refresh
// Method Post
func Refresh(w http.ResponseWriter, r *http.Request) {}
