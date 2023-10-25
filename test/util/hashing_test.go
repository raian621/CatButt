package util_test

import (
	"testing"

	"catbook.com/util"
)

func TestGenerateHash(t *testing.T) {
	_, err := util.DefaultGenerateHash("Hello")
	if err != nil {
		t.Error(err)
	}
}

func TestMatchHash(t *testing.T) {
	hash, err := util.DefaultGenerateHash("Lorem ipsum")
	if err != nil {
		t.Error(err)
	}
	matches, err := util.MatchHash("Lorem ipsum", hash)
	if err != nil {
		t.Error(err)
	}

	if !matches {
		t.Error("Hash was supposed to match original text")
	}

	matches, err = util.MatchHash("Not lorem upsum", hash)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(hash), "\n\n")

	if matches {
		t.Error("Hash was not supposed to match original text")
	}
}
