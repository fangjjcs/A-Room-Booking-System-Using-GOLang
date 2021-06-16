package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T){

	var myHandler myHandler

	h := NoSurf(&myHandler)
	switch v := h.(type){
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T",v))
	}
}

func TestSessionLoad(t *testing.T){
	var myHandler myHandler

	h := SessionLoad(&myHandler)
	switch result := h.(type){
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T",result))
	}
}