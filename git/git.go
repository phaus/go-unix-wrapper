package git

import (
	"errors"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/phaus/go-unix-wrapper/sys"
)

// Repository - a GIT Repository
type Repository struct {
	Name      string
	URL       string
	Folder    string
	LocalCopy string
}

// Bootstrap Bootstraps a local Copy.
func Bootstrap(argName string, argURL string, argFolder string) (Repository, error) {
	repo := Repository{Name: argName, URL: argURL, Folder: argFolder}
	localCopy, err := createLocalCopy(repo)
	if err != nil {
		return repo, err
	}
	repo.LocalCopy = localCopy
	git, err := sys.GetPath("git")
	if err != nil {
		return repo, err
	}
	cmd := exec.Command(git, "clone", repo.URL, repo.LocalCopy)
	out, err := sys.RunCmd(cmd)
	if err != nil {
		return repo, err
	}
	if out != "" {
		log.Printf("%s\n", out)
	}
	return repo, nil
}

// AddFile Adds a File to the Git Repository.
func AddFile(repo Repository, file string) (string, error) {
	log.Println("addFile")

	git, err := sys.GetPath("git")
	if err != nil {
		return "", err
	}

	addCmd := exec.Command(git, "add", file)
	addCmd.Dir = repo.LocalCopy
	addResult, addErr := sys.RunCmd(addCmd)
	if addErr != nil {
		return "", addErr
	}
	return addResult, nil
}

// CommitBranch Commits a Branch to the Git Repository
func CommitBranch(repo Repository, comment string) (string, error) {
	log.Println("commitBranch")

	git, err := sys.GetPath("git")
	if err != nil {
		return "", nil
	}

	cmtCmd := exec.Command(git, "commit", "-m", comment)
	cmtCmd.Dir = repo.LocalCopy
	cmtResult, cmtErr := sys.RunCmd(cmtCmd)
	if cmtErr != nil {
		return "", cmtErr
	}
	return cmtResult, nil
}

// CreateBranch creates a new Branch within the local Copy.
func CreateBranch(repo Repository, branch string) (string, error) {
	log.Printf("CreateBranch %s\n", branch)
	git, err := sys.GetPath("git")
	if err != nil {
		return "", err
	}
	_, err = pullRemote(repo, branch)
	var cmd *exec.Cmd
	if err != nil {
		cmd = exec.Command(git, "checkout", "-b", branch)
	} else {
		cmd = exec.Command(git, "checkout", branch)
	}
	cmd.Dir = repo.LocalCopy
	checkoutResult, checkoutErr := sys.RunCmd(cmd)
	if checkoutErr != nil {
		return "", checkoutErr
	}
	return checkoutResult, nil
}

// PushBranch pushes the changes of that branch to the remote Repository.
func PushBranch(repo Repository, branch string) (string, error) {
	log.Printf("PushBranch %s\n", branch)
	git, err := sys.GetPath("git")
	if err != nil {
		return "", err
	}
	pushCmd := exec.Command(git, "push", "-u", "origin", branch)
	pushCmd.Dir = repo.LocalCopy
	pushResult, pushErr := sys.RunCmd(pushCmd)
	if pushErr != nil {
		return "", pushErr
	}
	return pushResult, nil
}

// PullBranch pulls the changes of that branch from the remote Repository.
func PullBranch(repo Repository, branch string) (string, error) {
	log.Printf("PullBranch %s\n", branch)
	pullResult, pullErr := pullRemote(repo, branch)
	if pullErr != nil {
		return "", pullErr
	}
	return pullResult, nil
}

// Cleanup - removes a localCopy of a Repository.
func Cleanup(repo Repository) (string, error) {
	cmd := exec.Command("rm", "-Rf", repo.LocalCopy)
	out, err := sys.RunCmd(cmd)
	if err != nil {
		log.Fatalf("%s", err)
		return "", err
	}
	return out, nil
}

func pullRemote(repo Repository, branch string) (string, error) {
	git, err := sys.GetPath("git")
	if err != nil {
		return "", err
	}
	pullCmd := exec.Command(git, "pull", "origin", branch)
	pullCmd.Dir = repo.LocalCopy
	pullResult, pullErr := sys.RunCmd(pullCmd)
	if pullErr != nil {
		return "", err
	}
	return pullResult, nil
}

func createLocalCopy(repo Repository) (string, error) {
	if repo.Folder == "" {
		err := errors.New("folder must be set")
		log.Fatalf("%s", err)
		return "", err
	}
	cmd := exec.Command("mkdir", "-p", repo.Folder)
	out, err := sys.RunCmd(cmd)
	if err != nil {
		log.Fatalf("%s", err)
		return "", err
	}
	log.Printf("%s", out)
	localCopy, _ := filepath.Abs(repo.Folder)
	log.Printf("Crated localCopy: %s", localCopy)
	return localCopy, nil
}
