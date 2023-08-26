package main

import (
  "testing"
)


// "a4bc2d5e" => "aaaabccddddde"
// "abcd" => "abcd"
// "45" => "" (некорректная строка)
// "" => ""

func TestTask(t *testing.T) {
  var tests = []struct {
    input string
    expected string
  }{
    {"a4bc2d5e", "aaaabccddddde"},
    {"abcd", "abcd"},
    {"45", ""},
    {"", ""},
  }

  for _, test := range tests {
    if output := unpack(test.input); output != test.expected {
      t.Error("Test Failed: {} inputted, {} expected, received: {}", test.input, test.expected, output)
    }
  }
}
