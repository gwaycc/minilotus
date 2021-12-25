package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gwaylib/errors"
	"github.com/mitchellh/go-homedir"
)

func ExpandPath(repoDir string) string {
	r, err := homedir.Expand(repoDir)
	if err != nil {
		panic(err)
	}
	return r
}

const (
	REPO_TOKEN_FILE = "token"
)

type Repo struct {
	root string
}

func NewRepo(root string) (*Repo, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, errors.As(err)
	}
	return &Repo{root}, nil
}

func (repo *Repo) ReadToken() (string, error) {
	file := filepath.Join(repo.root, REPO_TOKEN_FILE)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", errors.As(err)
		}
		data := []byte(uuid.New().String())
		if err := ioutil.WriteFile(file, data, 0600); err != nil {
			return "", errors.As(err)
		}
	}
	return string(data), nil
}
