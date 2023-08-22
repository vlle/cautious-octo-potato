package main

import (
	"os/exec"
	"testing"
)

func originalGrep(pattern string, flags ...string) ([]byte, error) {
	app := "grep"
	flags = append(flags, pattern)
	cmd := exec.Command(app, flags...)
	stdout, _ := cmd.Output()
	return stdout, nil
}

func TestStuff(t *testing.T) {
	expected, err := originalGrep("example")
	my_result, my_err := grep("example")
	if err != my_err {
		t.Errorf("Err at TestStuff: %s", err.Error())
		t.Log(err)
		t.Log(my_err)
	}
	if len(expected) != len(my_result) {
		t.Errorf("expected != my_result")
		t.Log(expected)
		t.Log(my_result)
	}
}
