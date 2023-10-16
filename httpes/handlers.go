package httpes

import (
	"fmt"
	"net/http"
	"regexp"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func NewRegexpSolver() *regexpSolver {
	return &regexpSolver{
		handlers: make(map[string]http.HandlerFunc),
		cache:    make(map[string]*regexp.Regexp),
	}
}

type regexpSolver struct {
	handlers map[string]http.HandlerFunc
	cache    map[string]*regexp.Regexp
}

// Add
// add a new tamplate for a route on the server
func (s *regexpSolver) Add(regPath string, handler http.HandlerFunc) error {
	reg, err := regexp.Compile(regPath)
	if err != nil {
		return fmt.Errorf("registering regexp handler %s: %w", regPath, err)
	}
	s.handlers[regPath] = handler
	s.cache[regPath] = reg
	return nil
}

func (s *regexpSolver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	check := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	for pattern, handlerFunc := range s.handlers {
		if s.cache[pattern].MatchString(check) {
			handlerFunc(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// MapMid
// apply a list of middleware to a handlerfunc
// the middleware are applied on reverse order
func MapMid(next http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return next
	}
	warpedFunc := next
	//loopp on reverse order to preserve middleare order of exec
	for i := len(m) - 1; i >= 0; i-- {
		warpedFunc = m[i](warpedFunc)
	}
	return warpedFunc
}
