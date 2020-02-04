package git

import (
	"bytes"
	"errors"
	"fmt"
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
func Bootstrap(argName string, argURL string, argFolder string) (*Repository, error) {
	repo := &Repository{
		Name:   argName,
		URL:    argURL,
		Folder: argFolder}
	localCopy, err := repo.createLocalCopy()
	if err != nil {
		return repo, err
	}
	repo.LocalCopy = localCopy
	repo.Cleanup()
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
		log.Printf("%s", out)
	}
	return repo, nil
}

// AddFile Adds a File to the Git Repository.
func (repo *Repository) AddFile(file string) (string, error) {
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
func (repo *Repository) CommitBranch(comment string) (string, error) {
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
func (repo *Repository) CreateBranch(branch string) (string, error) {
	var buffer bytes.Buffer
	var out string
	log.Printf("CreateBranch %s\n", branch)
	git, err := sys.GetPath("git")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(git, "checkout", "-b", branch)
	cmd.Dir = repo.LocalCopy
	out, err = sys.RunCmd(cmd)
	if out != "" {
		buffer.WriteString(fmt.Sprintf("%s", out))
	}
	out, err = repo.resetLocalCopy(branch)
	if err != nil {
		return "", err
	}
	if out != "" {
		buffer.WriteString(fmt.Sprintf("%s", out))
	}
	return buffer.String(), nil
}

// PushBranch pushes the changes of that branch to the remote Repository.
func (repo *Repository) PushBranch(branch string) (string, error) {
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
func (repo *Repository) PullBranch(branch string) (string, error) {
	log.Printf("PullBranch %s\n", branch)
	pullResult, pullErr := repo.pullRemote(branch)
	if pullErr != nil {
		return "", pullErr
	}
	return pullResult, nil
}

// Cleanup - removes a localCopy of a Repository.
func (repo *Repository) Cleanup() (string, error) {
	cmd := exec.Command("rm", "-Rf", repo.LocalCopy)
	out, err := sys.RunCmd(cmd)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (repo *Repository) resetLocalCopy(branch string) (string, error) {
	var out string
	git, err := sys.GetPath("git")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(git, "fetch", "origin")
	cmd.Dir = repo.LocalCopy
	out, err = sys.RunCmd(cmd)
	if err != nil {
		return "", err
	}
	cmd = exec.Command(git, "reset", "--hard", fmt.Sprintf("origin/%s", branch))
	cmd.Dir = repo.LocalCopy
	out, err = sys.RunCmd(cmd)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (repo *Repository) pullRemote(branch string) (string, error) {
	git, err := sys.GetPath("git")
	if err != nil {
		return "", err
	}
	pullCmd := exec.Command(git, "pull", "origin", branch)
	pullCmd.Dir = repo.LocalCopy
	pullResult, pullErr := sys.RunCmd(pullCmd)
	if pullErr != nil {
		return "", pullErr
	}
	return pullResult, nil
}

func (repo *Repository) createLocalCopy() (string, error) {
	if repo.Folder == "" {
		err := errors.New("folder must be set")
		log.Printf("%s", err)
		return "", err
	}
	cmd := exec.Command("mkdir", "-p", repo.Folder)
	out, err := sys.RunCmd(cmd)
	if out != "" {
		log.Printf("%s", out)
	}
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}
	localCopy, _ := filepath.Abs(repo.Folder)
	log.Printf("Crated localCopy: %s", localCopy)
	return localCopy, nil
}
