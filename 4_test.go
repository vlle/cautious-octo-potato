package main

import (
	"testing"
)

func assertOk(t *testing.T, r map[string][]string, key string, l int) {
  if len(r[key]) != l {
    t.Errorf("Expected: len(r[\"%s\"] == 3, got: %d)", key, len(r[key]))
  }
}

func TestUnique(t *testing.T) {
  words := []string{"иван","вина","нива","рука", "лицо", "олиц", "кура"}
  r := f4(words)
  assertOk(t, r, "иван", 3)
}

func TestRegister(t *testing.T) {
  words := []string{"ИВАН","вина","НИВА","рука", "лицо", "олиц", "КУРА"}
  r := f4(words)
  assertOk(t, r, "иван", 3)
  assertOk(t, r, "рука", 2)
}

func TestDuplicates(t *testing.T) {
  words := []string{"ИВАН","иван", "ивАн"}
  r := f4(words)
  assertOk(t, r, "иван", 1)
}
