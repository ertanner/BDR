package main

import "testing"


func TestHello(t *testing.T){
	expected := "Hello Go!"
	actual := hello()
	if actual != expected {
		t.Error("Test failed")
	}
}

