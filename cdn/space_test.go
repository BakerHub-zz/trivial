package cdn

import (
	"github.com/BakerHub/trivial/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"testing"
)

type MockUploader struct {
	mock.Mock
}

func (m *MockUploader) Upload(path string, prefix string) []error {
	m.Called(path, prefix)
	return nil
}

func TestNewSpace(t *testing.T) {
	var u MockUploader
	space := NewSpace(&u, "unit/test/")
	assert.Equal(t, "unit/test/", space.prefix)
}

func TestSpace_Stage_success(t *testing.T) {
	cases := map[string]struct {
		prefix string
		stage  fs.MockDirectory
		hash   Hash
		file   *fs.MockFile

		cndPath string
	}{
		"sha1/empty file": {
			"prefix",
			".aspaces",
			fs.NewFileHashSHA1(),
			fs.NewMockFile("af.txt", ""),
			"prefix/da/39a3ee5e6b4b0d3255bfef95601890afd80709.txt",
		},
		"sha1/not ext": {
			"prefix/",
			".bspaces",
			fs.NewFileHashSHA1(),
			fs.NewMockFile("af", ""),
			"prefix/da/39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		"sha1/hidden file": {
			"prefix",
			".bspaces",
			fs.NewFileHashSHA1(),
			fs.NewMockFile(".bf", ""),
			"prefix/da/39a3ee5e6b4b0d3255bfef95601890afd80709.bf",
		},
		"normal file/sha1": {
			"prefix",
			".bspaces",
			fs.NewFileHashSHA1(),
			fs.NewMockFile("cf", "normal"),
			"prefix/9c/2a6e4809aeef7b7712ca4db05a681452f4f748",
		},
		"normal file/md5": {
			"prefix",
			".bspaces",
			fs.NewFileHashMD5(),
			fs.NewMockFile("cf", "normal"),
			"prefix/fe/a087517c26fadd409bd4b9dc642555",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			defer tc.stage.MustRemove(t)

			var u MockUploader
			s := &Space{&u, tc.prefix, tc.stage.Pathname(), tc.hash, map[string]string{}}

			tc.file.MustCreate(t)
			defer tc.file.MustRemove(t)

			cdnPath, err := s.Stage(tc.file.Pathname())
			assert.NoError(t, err)
			assert.Equal(t, tc.cndPath, cdnPath)

			staged := filepath.Join(tc.stage.Pathname(), cdnPath)
			assert.True(t, fs.Exists(staged), "expected file", staged)
		})
	}
}

func TestSpace_Stage_error_not_exits(t *testing.T) {
	notFile := "the-file-must-not-exit.txt"

	if fs.Exists(notFile) {
		err := os.Remove(notFile)
		assert.NoError(t, err, "error while remove file", notFile)
	}

	var u MockUploader
	s := NewSpace(&u, "test/prefix")
	_, err := s.Stage(notFile)

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestWithHash(t *testing.T) {
	h := fs.NewFileHashMD5()
	var u MockUploader
	s := NewSpace(&u, "test/prefix", WithHash(h))
	assert.Equal(t, h, s.hash)
}

func TestStageDirectory(t *testing.T) {
	var u MockUploader
	want := ".space-for-test"
	s := NewSpace(&u, "test/prefix", StageDirectory(want))
	assert.Equal(t, want, s.stageDirectory)
}

func TestSpace_Push(t *testing.T) {
	const (
		directory = "testdir"
		prefix    = "testprefix"
		path      = "testdir/testprefix"
	)
	m := &MockUploader{}
	s := NewSpace(m, prefix, StageDirectory(directory))
	m.On("Upload", path, prefix).Return(nil)
	s.Push()
	m.AssertExpectations(t)
}
