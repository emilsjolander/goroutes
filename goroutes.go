package goroutes

import (
  "net/http"
  "strings"
  "fmt"
)

var (
  // all the mapped routes
  routes []*route

  // otherwise know as the 404-function
  ResourceNotFoundHandler = func(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "Could not find this goroute!")
  }
)

// match a handler to a method and a pattern.
// this functions much like net/http packages Handle function
// diference is that this will match to a http method as well as allowing variables and wildcards in the pattern
// variables must be in their own path segment prefixed with a ':' ex: /user/:Id/settings 
// variables are extracted from Request.Form in the handlers ServeHTTP method. ParseForm must be called first. 
// patterns may define the last segment as a wildcard ex: /user/* which will match any path starting with /user/
// the method my have a zero value ("") which will match the pattern for any method
func Match(method string, pattern string, handler http.Handler) error {
  r, err := createRoute(method, pattern, handler)
  if err != nil {
    return err
  }
  routes = append(routes, r)
  return nil
}

// same as Match but takes in a function instead of an interface
func MatchFunc(method string, pattern string, handler func(http.ResponseWriter, *http.Request)) error {
  return Match(method, pattern, http.HandlerFunc(handler))
}

// this function will create RESTful url resources for a given controller with a set of parent controllers
// this function will call match with the following pattern and http method pairs.
// each pair will call a method of the controller (listed after the method-pattern pair)
// the example below is for a controller named ExampleController. 
// note that if the controller is not suffixed with "Controller" an error will be returned
//
// GET      /example             Index    
// GET      /example/new         New
// POST     /example             Create
// GET      /example/:Id         Show
// GET      /example/:Id/edit    Edit
// PUT      /example/:Id         Update
// DELETE   /example/:Id         Destroy
//
// if any parent controllers are given the preceding urls will the prefixed with /parent/:ParentId 
// you may specify any number of parent controllers.
// parent controllers should be string names of the actual controller name e.g. ParentController
func Resources(controller Controller, parentControllers ...string) error {

  parentPath, err := buildParentPath(parentControllers)
  if err != nil {
    return err
  }

  resourceName, err := getResourceName(getControllerName(controller))
  if err != nil {
    return err
  }
  controllerPath := parentPath + strings.ToLower(resourceName)

  MatchFunc("GET",    controllerPath + "/",         func(w http.ResponseWriter, r *http.Request){controller.Index(w,r)})
  MatchFunc("GET",    controllerPath + "/new",      func(w http.ResponseWriter, r *http.Request){controller.New(w,r)})
  MatchFunc("POST",   controllerPath + "/",         func(w http.ResponseWriter, r *http.Request){controller.Create(w,r)})
  MatchFunc("GET",    controllerPath + "/:Id",      func(w http.ResponseWriter, r *http.Request){controller.Show(w,r)})
  MatchFunc("GET",    controllerPath + "/:Id/edit", func(w http.ResponseWriter, r *http.Request){controller.Edit(w,r)})
  MatchFunc("PUT",    controllerPath + "/:Id",      func(w http.ResponseWriter, r *http.Request){controller.Update(w,r)})
  MatchFunc("DELETE", controllerPath + "/:Id",      func(w http.ResponseWriter, r *http.Request){controller.Destroy(w,r)})

  return nil
}

// build the path that precedes the controller name
func buildParentPath(parentControllers []string) (string, error) {
  path := "/"
  for i := len(parentControllers)-1; i>=0; i-- {
    c := parentControllers[i]
    name, err := getResourceName(c)
    if err != nil {
      return "",err
    }
    path += strings.ToLower(name) + "/:" + name+"Id/"
  }
  return path, nil
}

// this will try to match a path to a route
func handleRequest(w http.ResponseWriter, req *http.Request) {
  for _, r := range routes {
    if r.matches(req.Method, req.URL.Path) {
      r.handleRequest(w, req)
      return
    }
  }
  ResourceNotFound(w, req)
}

// call this to send a 404
// ovveride the behaviour of a 404 by setting ResourceNotFoundHandler to your own function
func ResourceNotFound(w http.ResponseWriter, req *http.Request) {
  ResourceNotFoundHandler(w, req)
}

// start listening for incoming trafic. this is a blocking operation
func Start() {
  http.HandleFunc("/", handleRequest)
  http.ListenAndServe(":9999", nil)
}
