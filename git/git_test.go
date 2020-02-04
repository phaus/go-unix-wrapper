package git_test

import (
	"testing"

	"github.com/phaus/go-unix-wrapper/git"
)

func TestBootstrap(t *testing.T) {
	repo, err := git.Bootstrap("go-unix-wrapper", "https://github.com/phaus/go-unix-wrapper.git", "data/go-unix-wrapper")
	if err != nil {
		t.Fatalf("%s", err)
	}
	repo.Cleanup()
}

func TestChangeBranch(t *testing.T) {
	repo, err := git.Bootstrap("go-unix-wrapper", "https://github.com/phaus/go-unix-wrapper.git", "data/go-unix-wrapper")
	if err != nil {
		t.Fatalf("%s", err)
	}
	_, err = repo.CreateBranch("master")
	if err != nil {
		t.Fatalf("while checking out master %s", err)
	}
	_, err = repo.CreateBranch("test")
	if err != nil {
		t.Fatalf("while checking out test %s", err)
	}
	repo.Cleanup()
}
