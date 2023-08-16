package main

import (
  "testing"
)

func TestF2_1(t *testing.T) {
    got := f2("a4bc2d5e")
    if got != "aaaabccddddde" {
        t.Errorf("f2(a4bc2d5e) = %s; want aaaabccddddde", got)
    }
}

func TestF2_2(t *testing.T) {
    got := f2("abcd")
    if got != "abcd"  {
        t.Errorf("f2(abcd) = %s; want abcd", got)
    }
}

func TestF2_3(t *testing.T) {
    got := f2("45")
    if got != ""  {
        t.Errorf("f2(45) = %s; want ", got)
    }
}

func TestF2_4(t *testing.T) {
    got := f2("")
    if got != ""  {
        t.Errorf("f2() = %s; want ", got)
    }
}
