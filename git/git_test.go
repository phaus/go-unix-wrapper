package git_test

import (
	"testing"

	"github.com/phaus/go-unix-wrapper/git"
)

func TestBootstrap(t *testing.T) {
	repo := git.Bootstrap("go-unix-wrapper", "https://github.com/phaus/go-unix-wrapper.git", "data/go-unix-wrapper")
	git.Cleanup(repo)
}
