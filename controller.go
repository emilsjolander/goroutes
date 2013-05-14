package goroutes

import (
  "net/http"
  "strings"
  "regexp"
  "errors"
  "reflect"
  "bitbucket.org/pkg/inflect"
)

// All possible controller actions
type Action uint
const (
  Index Action = iota
  New
  Create
  Show
  Edit
  Update
  Destroy
)

func (a Action) String() string {
  if int(a) < len(actionNames) {
    return actionNames[a]
  }
  return ""
}

var actionNames = []string{
  Index:        "Index",
  New:          "New",
  Create:       "Create",
  Show:         "Show",
  Edit:         "Edit",
  Update:       "Update",
  Destroy:      "Destroy",
}

var (
  //validates the charachters alowed in a controller name. might be extended in the future
  controllerNameValidator = regexp.MustCompile("[a-zA-Z]+")
)

type BeforeFilterer interface{
  BeforeFilter(a Action, w http.ResponseWriter, r *http.Request) bool
}

type Controller interface {
  Index   (w http.ResponseWriter, req *http.Request)
  New     (w http.ResponseWriter, req *http.Request)
  Create  (w http.ResponseWriter, req *http.Request)
  Show    (w http.ResponseWriter, req *http.Request)
  Edit    (w http.ResponseWriter, req *http.Request)
  Update  (w http.ResponseWriter, req *http.Request)
  Destroy (w http.ResponseWriter, req *http.Request)
}

// a base implementation of a controller that returns a 404 for every method
// this is usefull for having as an anonymous field in your controller
// that way you don't need to implement the methods that you do not use.
type BaseController struct {}

func (c *BaseController) Index(w http.ResponseWriter, req *http.Request) {
  ResourceNotFound(w, req)
}

func (c *BaseController) New(w http.ResponseWriter, req *http.Request) {
  ResourceNotFound(w, req)
}

func (c *BaseController) Create(w http.ResponseWriter, req *http.Request) {
  ResourceNotFound(w, req)
}

func (c *BaseController) Show(w http.ResponseWriter, req *http.Request) {
  ResourceNotFound(w, req)
}

func (c *BaseController) Edit(w http.ResponseWriter, req *http.Request) {
  ResourceNotFound(w, req)
}

func (c *BaseController) Update(w http.ResponseWriter, req *http.Request) {
  ResourceNotFound(w, req)
}

func (c *BaseController) Destroy(w http.ResponseWriter, req *http.Request) {
  ResourceNotFound(w, req)
}

// reflects on the controller to get the name of the underlying struct
func getControllerName(controller Controller) string {
  name := ""
  if reflect.TypeOf(controller).Kind() == reflect.Ptr {
    name = reflect.TypeOf(controller).Elem().Name()
  }else {
    name = reflect.TypeOf(controller).Name()
  }
  return name
}

// gets a resource name from a controller name, returning an error if the controller name is invalid.
// an example input -> output is "ExampleController" -> "example"
// an error would occur from input "Example"
func getResourceName(controllerName string) (string, error) {
  if !strings.HasSuffix(controllerName, "Controller") {
    return "", errors.New("goroutes: controller name must have suffix \"Controller\"")
  }
  controllerName = controllerName[:len(controllerName)-len("Controller")]
  controllerNameUnderscore := inflect.Underscore(controllerName)
  if controllerNameUnderscore == "new" || controllerNameUnderscore == "edit" || !controllerNameValidator.MatchString(controllerName) {
    return "", errors.New("goroutes: controller name must not be NewController, EditController or include invalid characters")
  }
  return controllerNameUnderscore,nil 
}
