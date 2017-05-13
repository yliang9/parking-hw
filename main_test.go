package main

import (
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Init()
	code := m.Run()
	os.Exit(code)
}
