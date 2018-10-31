package git

import (
	"log"
	"os/exec"
	"path/filepath"

	"gitlab.innoq.com/phl/forwarder/sys"
)

var localCopy string
var repositoryURL string

// Bootstrap Bootstraps a local Copy.
func Bootstrap(repoURL string) {
	repositoryURL = repoURL
	localCopy, _ = createLocalCopy()
	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}
	cmd := exec.Command(git, "clone", repositoryURL, localCopy)
	out, err := sys.RunCmd(cmd)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n", out)
}

// AddFile Adds a File to the Git Repository.
func AddFile(file string) {
	log.Println("addFile")

	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}

	addCmd := exec.Command(git, "add", file)
	addCmd.Dir = localCopy
	addResult, addErr := sys.RunCmd(addCmd)
	if addErr != nil {
		log.Println(addErr)
	}
	log.Printf("%s\n", addResult)

}

// CommitBranch Commits a Branch to the Git Repository
func CommitBranch(comment string) {
	log.Println("commitBranch")

	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}

	cmtCmd := exec.Command(git, "commit", "-m", comment)
	cmtCmd.Dir = localCopy
	cmtResult, cmtErr := sys.RunCmd(cmtCmd)
	if cmtErr != nil {
		log.Println(cmtErr)
	}
	log.Printf("%s\n", cmtResult)
}

// CreateBranch creates a new Branch within the local Copy.
func CreateBranch(branch string) {
	log.Println("CreateBranch")
	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}
	_, err = pullRemote(branch)
	var cmd *exec.Cmd
	if err != nil {
		cmd = exec.Command(git, "checkout", "-b", branch)
	} else {
		cmd = exec.Command(git, "checkout", branch)
	}
	cmd.Dir = localCopy
	checkoutResult, checkoutErr := sys.RunCmd(cmd)
	if checkoutErr != nil {
		log.Println(checkoutErr)
	}
	log.Printf("%s\n", checkoutResult)
}

// PushBranch pushes the changes of that branch to the remote Repository.
func PushBranch(branch string) {
	log.Println("PushBranch")
	git, err := sys.GetPath("git")
	if err != nil {
		log.Println(err)
	}
	pushCmd := exec.Command(git, "push", "-u", "origin", branch)
	pushCmd.Dir = localCopy
	pushResult, pushErr := sys.RunCmd(pushCmd)
	if pushErr != nil {
		log.Println(pushErr)
	}
	log.Printf("%s\n", pushResult)
}

// PullBranch pulls the changes of that branch from the remote Repository.
func PullBranch(branch string) {
	log.Println("PullBranch")
	pullResult, pullErr := pullRemote(branch)
	if pullErr != nil {
		log.Println(pullErr)
	}
	log.Printf("%s\n", pullResult)
}

func pullRemote(branch string) (string, error) {
	git, err := sys.GetPath("git")
	if err != nil {
		return "", err
	}
	pullCmd := exec.Command(git, "pull", "origin", branch)
	pullCmd.Dir = localCopy
	pullResult, pullErr := sys.RunCmd(pullCmd)
	if pullErr != nil {
		return "", err
	}
	log.Printf("%s\n", pullResult)
	return pullResult, nil
}

// LocalCopy Return the path of the local Repository copy.
func LocalCopy() string {
	return localCopy
}

func createLocalCopy() (folder string, err error) {
	cmd := exec.Command("mkdir", "-p", "data")
	out, err := sys.RunCmd(cmd)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Printf("%s\n", out)
	return filepath.Abs("data")
}
