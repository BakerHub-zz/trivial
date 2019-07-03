package cdn

import (
	"github.com/BakerHub/trivial/shell"
	"strings"
)

type QShellUploader struct {
	bucket       string
	skipSuffixes string
}

func (qs *QShellUploader) Upload(directory string, prefix string) []error {
	shell.Run("qshell",
		"qupload2",
		"--src-dir", directory,
		"--bucket", qs.bucket,
		"--key-prefix", prefix,
		"--rescan-local",
		"--skip-suffixes", qs.skipSuffixes,
	)
	return nil
}

type QShellUploaderOption func(*QShellUploader)

func NewQShellUploader(bucket string, opts ...QShellUploaderOption) *QShellUploader {
	qs := &QShellUploader{bucket, ".DS_Store,Thumbs.db"}

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
