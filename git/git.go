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
func Bootstrap(repo Repository) Repository {
	localCopy, err := createLocalCopy(repo)
	if err != nil {
		log.Println(err)
	}
	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}
	cmd := exec.Command(git, "clone", repo.URL, repo.LocalCopy)
	out, err := sys.RunCmd(cmd)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n", out)
	return repo
}

// AddFile Adds a File to the Git Repository.
func AddFile(repo Repository, file string) {
	log.Println("addFile")

	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}

	addCmd := exec.Command(git, "add", file)
	addCmd.Dir = repo.LocalCopy
	addResult, addErr := sys.RunCmd(addCmd)
	if addErr != nil {
		log.Println(addErr)
	}
	log.Printf("%s\n", addResult)

}

// CommitBranch Commits a Branch to the Git Repository
func CommitBranch(repo Repository, comment string) {
	log.Println("commitBranch")

	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}

	cmtCmd := exec.Command(git, "commit", "-m", comment)
	cmtCmd.Dir = repo.LocalCopy
	cmtResult, cmtErr := sys.RunCmd(cmtCmd)
	if cmtErr != nil {
		log.Println(cmtErr)
	}
	log.Printf("%s\n", cmtResult)
}

// CreateBranch creates a new Branch within the local Copy.
func CreateBranch(repo Repository, branch string) {
	log.Printf("CreateBranch %s\n", branch)
	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
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
		log.Println(checkoutErr)
	}
	log.Printf("%s\n", checkoutResult)
}

// PushBranch pushes the changes of that branch to the remote Repository.
func PushBranch(repo Repository, branch string) {
	log.Printf("PushBranch %s\n", branch)
	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}
	pushCmd := exec.Command(git, "push", "-u", "origin", branch)
	pushCmd.Dir = repo.LocalCopy
	pushResult, pushErr := sys.RunCmd(pushCmd)
	if pushErr != nil {
		log.Println(pushErr)
	}
	log.Printf("%s\n", pushResult)
}

// PullBranch pulls the changes of that branch from the remote Repository.
func PullBranch(repo Repository, branch string) {
	log.Printf("PullBranch %s\n", branch)
	pullResult, pullErr := pullRemote(repo, branch)
	if pullErr != nil {
		log.Println(pullErr)
	}
	log.Printf("%s\n", pullResult)
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
	log.Printf("%s\n", pullResult)
	return pullResult, nil
}

func createLocalCopy(repo Repository) (string, error) {
	if repo.Folder == "" {
		err := errors.New("folder must be set")
		log.Fatal(err)
		return "", err
	}
	cmd := exec.Command("mkdir", "-p", repo.Folder)
	out, err := sys.RunCmd(cmd)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Printf("%s\n", out)
	return filepath.Abs(repo.Folder)
}
