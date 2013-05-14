package goroutes

import (
  "testing"
)

type SuperUserController struct {
  BaseController
}

func TestGetControllerName(t *testing.T){
  if getControllerName(new(SuperUserController)) != "SuperUserController" {
    t.Fail()
  }
}

func TestGetResourceName(t *testing.T){
  if n, _ := getResourceName("SuperUserController"); n != "super_user" {
    t.Fail()
  }
  if n, _ := getResourceName("UsersController"); n != "users" {
    t.Error(n)
  }
  if _, err := getResourceName("SuperUserCont"); err == nil {
    t.Fail()
  }
  if _, err := getResourceName("NewController"); err == nil {
    t.Fail()
  }
  if _, err := getResourceName("EditController"); err == nil {
    t.Fail()
  }
}

