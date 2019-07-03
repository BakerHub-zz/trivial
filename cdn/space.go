package cdn

import (
	"fmt"
	"github.com/BakerHub/trivial/fs"
	"os"
	"path/filepath"
)

type Hash interface {
	FromFile(path string) (string, error)
}

type Uploader interface {
	Upload(path string, prefix string) []error
}

type Space struct {
	uploader       Uploader
	prefix         string
	stageDirectory string
	hash           Hash

	files map[string]string
}

type Option func(space *Space)

func StageDirectory(stageDirectory string) Option {
	return func(space *Space) {
		space.stageDirectory = stageDirectory
	}
}

func WithHash(hash Hash) Option {
	return func(space *Space) {
		space.hash = hash
	}
}

// NewSpace creates a new cdn space that used to upload files.
func NewSpace(uploader Uploader, prefix string, options ...Option) *Space {
	space := &Space{
		uploader:       uploader,
		prefix:         prefix,
		stageDirectory: ".spaces",
		hash:           fs.NewFileHashSHA1(),
	}

	for _, option := range options {
		option(space)
	}

	return space
}

func (space *Space) hashFile(path string) (string, error) {
	return space.hash.FromFile(path)
}

func (space *Space) makeCdnPath(hash string, extension string) string {
	filename := fmt.Sprintf("%s%s", hash[2:], extension)
	return filepath.Join(space.prefix, hash[0:2], filename)
}

func (space *Space) cdnPath(localFile string) (string, error) {
	h, err := space.hashFile(localFile)
	if err != nil {
		return "", err
	}
	return space.makeCdnPath(h, filepath.Ext(localFile)), nil
}

func (space *Space) copyToStage(localFile, cndPath string) error {
	stagePath := filepath.Join(space.stageDirectory, cndPath)
	err := os.MkdirAll(filepath.Dir(stagePath), os.ModePerm)
	if err != nil {
		return err
	}
	return fs.CopyFile(localFile, stagePath)
}

// Stage put one local file to the stage area and return the final cnd path the file will finally upload to.
func (space *Space) Stage(localFile string) (cdnPath string, err error) {
	cdnPath, err = space.cdnPath(localFile)
	if err != nil {
		return "", err
	}

	err = space.copyToStage(localFile, cdnPath)
	if err != nil {
		return "", err
	}

	return cdnPath, nil
}

func (space *Space) stagedRoot() string {
	return filepath.Join(space.stageDirectory, space.prefix)
}

func (space *Space) Push() []error {
	return space.uploader.Upload(space.stagedRoot(), space.prefix)
}
