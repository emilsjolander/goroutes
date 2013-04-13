Goroutes
=========

About
-----

Goroutes is a url routing library for go with support for routing RESTful routes to controllers.

This started a with me wanting to experiment with go, but i have been developing using rails for the past months and a quickly missed a lot of its features. Url routing with variables and RESTful routes to controllers were one of the major features i missed so i decided to build this RESTful routing library for go.

It takes inspiration from rails as well as the standard go url routing found in net/http.


Usage
-----

First thing you have to do is get the project
```shell
go get github.com/emilsjolander/goroutes
```

After that all you need to do is import it into your project
```go
import "github.com/emilsjolander/goroutes"
```

Done! Now just match some routes and press Start!
```go
goroutes.Resources(new(UsersController))
goroutes.Resources(new(NotesController), "UsersController")

goroutes.MatchFunc("GET", "/status", 
  func(w http.ResponseWriter, req *http.Request){
    fmt.Fprintf(w, "Status ok!")
  })

goroutes.Match("GET", "assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
goroutes.Match("GET", "/", new(HomeHandler))

goroutes.Start()
```

Api
---

I have designed the api to be as similar as possible to go's routing while at the same time adding as many of the great routing features that come with a framework like rails.

There is only one public interface and that is the Controller interface. If your datastructure implements this interface goroutes can match RESTful routes the the corrosponding methods.
```go
type Controller interface {
  Index   (w http.ResponseWriter, req *http.Request)
  New     (w http.ResponseWriter, req *http.Request)
  Create  (w http.ResponseWriter, req *http.Request)
  Show    (w http.ResponseWriter, req *http.Request)
  Edit    (w http.ResponseWriter, req *http.Request)
  Update  (w http.ResponseWriter, req *http.Request)
  Destroy (w http.ResponseWriter, req *http.Request)
}
```

A lot of times a controller will only want to implement a handful of these methods, for this i have defined a BaseController struct that can be included into any controller as an anonymous field. Any method not overriden in your struct will answer the request with a 404.
```go
type MyController struct {
    goroutes.BaseController
}

func (c *MyController) Index(w http.ResponseWriter, req *http.Request) {
    // my response
}
```

To generate resources for a controller call
```go
goroutes.Resources(new(ExampleController))
```
This call will generate the following resources
```text
GET      /example             (Index)    
GET      /example/new         (New)
POST     /example             (Create)
GET      /example/:Id         (Show)
GET      /example/:Id/edit    (Edit)
PUT      /example/:Id         (Update)
DELETE   /example/:Id         (Destroy)
```

If any number of parent controllers were given that would prefix every url pattern with /parent/:ParentId, this can be achieved by the following method call.
```go
goroutes.Resources(new(ExampleController), "ParentController")
```

Multiple parent controllers are also supported
```go
goroutes.Resources(new(ExampleController), "ParentController", "GrandparentController")
```

The preceding call would generate the following resources
```text
GET      /grandparent/:GrandparentId/parent/:ParentId/example             (Index)    
GET      /grandparent/:GrandparentId/parent/:ParentId/example/new         (New)
POST     /grandparent/:GrandparentId/parent/:ParentId/example             (Create)
GET      /grandparent/:GrandparentId/parent/:ParentId/example/:Id         (Show)
GET      /grandparent/:GrandparentId/parent/:ParentId/example/:Id/edit    (Edit)
PUT      /grandparent/:GrandparentId/parent/:ParentId/example/:Id         (Update)
DELETE   /grandparent/:GrandparentId/parent/:ParentId/example/:Id         (Destroy)
```

Controllers are not always the correct solution so there are two more methods for routing urls.
They are really the same method only one takes in a struct implementing http.Handler and the other a handler function.
```go
func Match(method string, pattern string, handler http.Handler) error
func MatchFunc(method string, pattern string, handler func(http.ResponseWriter, *http.Request)) error 
```

They expept a http method (GET, POST, PUT or DELETE) or a empty string to indicate that the handler will handle all methods.
The pattern sent in can contain variables which are preceded with a ':'. The last segment or the url may also be the wildcard character '*'. The wilcard will match anything after it while the variables will only match anything in the corresponding segment.

Pattern variables are extracted from the Request the same way form and query params are extracted.
```go
func (c *ExampleController) Show(w http.ResponseWriter, req *http.Request) {
  req.ParseForm()
  id := req.Form["Id"][0]
  parentId := req.Form["ParentId"][0]
}
```


Contributing
------------

Pull requests and issues are very welcome!

Feature request are also welcome but i can't make any promise that they will make it in.
I would like to keep the library as general as possible, if you are unsure you can just ask before you code ;)


License
-------

    Copyright 2013 Emil Sjölander

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
