package main

import (
	"testing"
)

func BuildRunner() Runner {
	c := Runner{}
	c.Inject(Store{}, &InMemoryBackend{})
	return c
}

func TestBoomVersion(*testing.T) {
	runner := BuildRunner()
	err := runner.Delegate("version", "", "")

	if err {
		t.Fail("Error returned")
	}
}
