package git

import (
	"os/exec"
	"strings"
)

type VCS interface {
	Init() error
	Commit(msg string) error
	Push() (string, error)
	Pull() (string, error)
	Log() ([]string, error)
	DiffLast() (string, error)
}
type GitVCS struct{ File string }

func (g *GitVCS) Init() error {
	exec.Command("git", "init").Run()
	exec.Command("git", "add", ".").Run()
	return exec.Command("git", "commit", "-m", "Init").Run()
}
func (g *GitVCS) Commit(msg string) error {
	exec.Command("git", "add", ".").Run()
	return exec.Command("git", "commit", "-m", msg).Run()
}
func (g *GitVCS) Push() (string, error) {
	out, err := exec.Command("git", "push").CombinedOutput()
	return string(out), err
}
func (g *GitVCS) Pull() (string, error) {
	out, err := exec.Command("git", "pull").CombinedOutput()
	return string(out), err
}
func (g *GitVCS) Log() ([]string, error) {
	out, _ := exec.Command("git", "log", "-n", "10", "--pretty=format:%h - %s (%cr)").Output()
	if len(out) == 0 {
		return []string{}, nil
	}
	return strings.Split(string(out), "\n"), nil
}
func (g *GitVCS) DiffLast() (string, error) {
	out, err := exec.Command("git", "show", "--pretty=format:", "HEAD").CombinedOutput()
	return string(out), err
}
