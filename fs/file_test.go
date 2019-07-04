package fs

import (
	"errors"
	"github.com/BakerHub/trivial/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"os"
	"testing"
)

type MockedCloser struct {
	mock.Mock
}


func (m *MockedCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestMustCloseSuccess(t *testing.T) {
	h := new(MockedCloser)
	h.On("Close").Return(nil)

	MustClose(h)

	h.AssertExpectations(t)
}

func TestMustClosePanic(t *testing.T) {
	h := new(MockedCloser)
	h.On("Close").Return(errors.New("bad happened"))
	assert.Panics(t, func() {
		MustClose(h)
	})
	h.AssertExpectations(t)
}


func TestMustCopyFileSuccess(t *testing.T)  {
	src := &MockFile{"a.txt", "this this the content."}
	src.MustCreate(t)
	defer src.MustRemove(t)

	MustCopyFile(src.Pathname(), "any.txt")

	content, err := ioutil.ReadFile("any.txt")
	assert.NoError(t, err)
	assert.Equal(t, src.Content(), string(content))

	err = os.Remove("any.txt")
	assert.NoError(t, err)
}


func TestMustCopyFileFail(t *testing.T)  {
	assert.Panics(t, func() {
		MustCopyFile("not exists", "any")
	})
	assert.False(t, Exists("any"), "destination should not exist")
}

func TestCopyFileSuccess(t *testing.T)  {
	src := &MockFile{"a.txt", "this this the content."}
	src.MustCreate(t)
	defer src.MustRemove(t)

	err := CopyFile(src.Pathname(), "any.txt")
	assert.NoError(t, err)

	content, err := ioutil.ReadFile("any.txt")
	assert.NoError(t, err)
	assert.Equal(t, src.Content(), string(content))
	err = os.Remove("any.txt")
	assert.NoError(t, err)
}

func TestCopyFileFail(t *testing.T)  {
	err := CopyFile("not exists", "any")
	assert.Error(t, err)
	assert.False(t, Exists("any"), "destination should not exist")
}


func TestExists(t *testing.T) {
	cases := map[string]struct {
		path string
		want bool
	}{
		"exist file": {"file_test.go", true},
		"exist dir":  {".", true},
		"not exist":  {"file-not-exist.go", false},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := Exists(tc.path)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEnsureDirSuccess(t *testing.T) {
	dirs := []MockDirectory{MockDirectory("test-not-exits"), MockDirectory("test-exits")}
	dirs[0].MustCreate(t)
	dirs[1].MustCreate(t)

	for _, dir := range dirs {
		EnsureDir(dir.Pathname())
		assert.True(t, Exists(dir.Pathname()))
		dir.MustRemove(t)
	}
}

func TestListDirectory(t *testing.T)  {
	cases := map[string]struct{
		setup func()
		teardown func()
		root string
		want []string
	} {
		"directory": {
			func() {
				check.Check(os.MkdirAll("a/b/c", 0750))
				check.Check(ioutil.WriteFile("a/b.txt", []byte("hah"), 0777))
				check.Check(ioutil.WriteFile("a/b/c.txt", []byte("hah"), 0777))
			},
			func() {
				check.Check(os.RemoveAll("a/b/c"))
				check.Check(os.RemoveAll("a/b"))
				check.Check(os.RemoveAll("a"))
			},
				"a",
				[]string{
					"a/b.txt",
					"a/b/c.txt",
				},
		},
		"no files": {
			func() {
				check.Check(os.MkdirAll("a/b/c", 0750))
			},
			func() {
				check.Check(os.RemoveAll("a/b/c"))
				check.Check(os.RemoveAll("a/b"))
				check.Check(os.RemoveAll("a"))
			},
			"a",
			nil,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.setup()
			defer tc.teardown()

			got := ListDirectory(tc.root)
			assert.ElementsMatch(t, tc.want, got)
		})
	}
}

func TestReadLines(t *testing.T) {
	cases := map[string]struct{
		file *MockFile
		want []string
	} {
		"empty": {
			&MockFile{"a.txt", ""},
			nil,
		},
		"text lines": {
			&MockFile{"a.txt", "a\nb\nc"},
			[]string{"a","b", "c"},
		},
		"empty lines": {
			&MockFile{"a.txt", "a\n\nc"},
			[]string{"a","", "c"},
		},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc.file.MustCreate(t)
			defer tc.file.MustRemove(t)
			got, err := ReadLines(tc.file.Pathname())
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
