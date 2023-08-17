package main

import (
  "testing"
  "os/exec"
)

func originalGrep(pattern string, flags ...string) ([]byte, error) {
  app := "grep"
  flags = append(flags, pattern)
  cmd := exec.Command(app, flags...)
  stdout, err := cmd.Output()
  return stdout, err 
}

func TestStuff(t *testing.T) {
  expected, err := originalGrep("example")
  my_result, my_err := grep("example")
  if err != my_err {
    t.Errorf("Err at TestStuff: %s", err.Error())
  }
  if len(expected) != len(my_result) {
    t.Errorf("expected != my_result")
    t.Log(expected)
    t.Log(my_result)
  }
}
