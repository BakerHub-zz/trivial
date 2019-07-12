package cdn

import (
	"github.com/BakerHub/trivial/shell"
	"strings"
)

type QShellUploader struct {
	shell 			 Runner
	bucket       string
	skipSuffixes string
	local		 bool
}

type Runner interface {
	Run(name string, args...string)
}

func (qs *QShellUploader) Upload(directory string, prefix string) []error {
	args := []string {
		"qupload2",
		"--src-dir", directory,
		"--bucket", qs.bucket,
		"--key-prefix", prefix,
		"--rescan-local",
		"--skip-suffixes", qs.skipSuffixes,
	}
	if qs.local {
		args = append(args, "--local")
	}

	qs.shell.Run("qshell", args...)

	return nil
}

type QShellUploaderOption func(*QShellUploader)

func NewQShellUploader(bucket string, opts ...QShellUploaderOption) *QShellUploader {
	qs := &QShellUploader{
		shell: &shell.Shell{},
		bucket: bucket,
		skipSuffixes: ".DS_Store,Thumbs.db",
	}

	for _, opt := range opts {
		opt(qs)
	}

	return qs
}

func IgnoreSuffixes(suffixes ...string) QShellUploaderOption {
	all := append(suffixes, ".DS_Store", "Thumbs.db")
	ignores := strings.Join(all, ",")

	return func(qs *QShellUploader) {
		qs.skipSuffixes = ignores
	}
}

func Local() QShellUploaderOption {
	return func(qs *QShellUploader) {
		qs.local = true
	}
}
