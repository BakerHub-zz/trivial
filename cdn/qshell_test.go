package cdn

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestWithIgnoreSuffixes(t *testing.T) {
	cases := map[string]struct {
		args []string

		want string
	}{
		"empty": {
			[]string{},
			".DS_Store,Thumbs.db",
		},
		"append": {
			[]string{"a", "b"},
			"a,b,.DS_Store,Thumbs.db",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			var q QShellUploader
			opt := IgnoreSuffixes(tc.args...)
			opt(&q)
			assert.Equal(t, tc.want, q.skipSuffixes)
		})
	}
}

func TestNewQShellUploader(t *testing.T) {
	NewQShellUploader("test")
	NewQShellUploader("test", IgnoreSuffixes(".exe"))
}

type MockShell struct {
	mock.Mock
}

func (m *MockShell) Run(name string, runArgs...string) {
	m.Called(name, runArgs)
}

func WithMockShell(m * MockShell) QShellUploaderOption  {
	return func(uploader *QShellUploader) {
		uploader.shell = m
	}
}

type ShellCommand struct {
	name string
	args []string
}
func TestQShellUploader_Upload(t *testing.T) {
	cases := map[string]struct{
		bucket string
		options []QShellUploaderOption
		directory string
		prefix string
		want *ShellCommand
	} {
		"default options": {
			"bucket",
			nil,
			"test",
			"prefix",
			&ShellCommand{
				"qshell",
				[]string{"qupload2", "--src-dir", "test", "--bucket", "bucket", "--key-prefix", "prefix", "--rescan-local", "--skip-suffixes", ".DS_Store,Thumbs.db"},
			},
		},
		"with local options": {
			"bucket",
			[]QShellUploaderOption{Local()},
			"test",
			"prefix",
			&ShellCommand{
				"qshell",
				[]string{"qupload2", "--src-dir", "test", "--bucket", "bucket", "--key-prefix", "prefix", "--rescan-local", "--skip-suffixes", ".DS_Store,Thumbs.db", "--local"},
			},
		},
		"with ignore suffix": {
			"bucket",
			[]QShellUploaderOption{IgnoreSuffixes("abc", "def")},
			"test",
			"prefix",
			&ShellCommand{
				"qshell",
				[]string{"qupload2", "--src-dir", "test", "--bucket", "bucket", "--key-prefix", "prefix", "--rescan-local", "--skip-suffixes", "abc,def,.DS_Store,Thumbs.db"},
			},
		},
	}

	for name, tc := range cases  {
		tc := tc
		t.Run(name, func(t *testing.T) {
			s := MockShell{}
			options := append(tc.options, WithMockShell(&s))
			qs := NewQShellUploader(tc.bucket, options...)
			s.On("Run", tc.want.name, tc.want.args).Return()

			qs.Upload(tc.directory, tc.prefix)

			s.AssertExpectations(t)
		})
	}
}
