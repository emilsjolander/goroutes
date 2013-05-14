package goroutes

import (
  "net/http"
  "strings"
  "regexp"
  "errors"
)

var (
  // validates that a pattern contains valid charachters
  patternValidator  = regexp.MustCompile("^[\\w\\d-\\/:\\*]+$")
  // used when replacing a pattern variable
  varReplacer       = regexp.MustCompile("\\\\/:[\\w\\d-]+\\\\/")
  // used when replacing a pattern wildcard
  wildcardReplacer  = regexp.MustCompile("\\\\/\\*\\\\/")
  // used to escape slashen in a pattern
  slashEscaper      = regexp.MustCompile("/")
)

type route struct {
  // GET PUT POST or DELETE. May also be a empty string to match any method
  method          string
  // the raw pattern
  pattern         string
  // the regexp used to check if this route should handle a request
  patternMatcher  *regexp.Regexp
  //the handler responsible for handling requests that match this route
  handler         http.Handler
}

// creates a new route. will instantiate a bunch for regexp and 
// return an error if the pattern is not valid
func createRoute(method string, pattern string, handler http.Handler) (*route, error) {

  if pattern[len(pattern)-1] != '/' {
    pattern += "/"
  }

  err := validatePattern(pattern)
  if err != nil {
    return nil, err
  }

  regexpPattern := strings.ToLower(pattern)
  regexpPattern  = slashEscaper.ReplaceAllString(regexpPattern, "\\/")
  regexpPattern  = varReplacer.ReplaceAllString(regexpPattern, "\\/[\\w\\d-]+\\/")
  regexpPattern  = wildcardReplacer.ReplaceAllString(regexpPattern, "\\/[\\w\\d-\\/\\.]+\\/")
  regexpPattern  = "^" + regexpPattern + "?$"

  r := &route{  method:         method,
                pattern:        pattern,
                handler:        handler,
                patternMatcher: regexp.MustCompile(regexpPattern) }

  return r, nil
}

// check so a pattern is valid
// this includes checking so it only contains certain characters
// also check to variables are listed in their own segment
// and so wildcards are alone in the last segment
func validatePattern(pattern string) error {
  if !patternValidator.MatchString(pattern) {
    return errors.New("goroutes: pattern contains invalid characters")
  }
  for i := 0; i<len(pattern); i++ {
    if pattern[i] == ':' && (i == 0 || pattern[i-1] != '/') {
      return errors.New("goroutes: a variable was not alone in it's segment")
    }else if pattern[i] == '*' && (i != len(pattern)-2 || pattern[i-1] != '/') {
      return errors.New("goroutes: wildcard must be the last segment")
    }
  }
  return nil
}

// check the the url matches this route
func (r *route) matches(method string, url string) bool {
  return (r.method == "" || method == r.method) && r.patternMatcher.MatchString(url)
}

// handle the request, this inclused extracting the pattern variables from the url
// the pattern variables will be contatinated with the query making them available after
// ParseForm has been called in the request
func (r *route) handleRequest(w http.ResponseWriter, req *http.Request) {
  pathParams := r.parsePatternParams(req.URL.Path)
  if req.URL.RawQuery != "" && pathParams != "" {
    req.URL.RawQuery += "&"
  }
  req.URL.RawQuery += pathParams
  r.handler.ServeHTTP(w,req)
}

// parser the url and extracts pattern variables and their value
func (r *route) parsePatternParams(path string) string {
  pathParts     := strings.Split(path, "/")
  patternParts  := strings.Split(r.pattern, "/")
  params := ""
  for i, s := range patternParts {
    if strings.HasPrefix(s, ":") {
      if params != "" {
        params += "&"
      }
      params += s[1:] + "=" + pathParts[i]
    }
  }
  return params
}


