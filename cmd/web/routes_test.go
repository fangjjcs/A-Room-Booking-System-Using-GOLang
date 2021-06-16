package main

import (
	"fmt"
	"testing"

	"github.com/fangjjcs/bookings-app/pkg/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T){
	var app config.AppConfig
	mux := routes(&app)

	switch result := mux.(type){
	case *chi.Mux:
		// do nothing
	default:
		t.Error(fmt.Sprintf("Type is not *chi.Mux, but it's %T", result))
	}
}