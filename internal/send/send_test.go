package send

import (
	"os"
	"path/filepath"
	"testing"
)

type TestValidPathCase struct {
	path   string
	result bool
	err    error
}

func TestValidPath(t *testing.T) {
	testCase := []TestValidPathCase{
		{
			result: true,
			path:   "x.mod",
			err:    nil,
		},
		{
			result: true,
			path:   "y.sum",
			err:    nil,
		},
		{
			result: true,
			path:   "z.go",
			err:    nil,
		},
		{
			result: false,
			path:   "xxxx.go",
			err:    os.ErrNotExist,
		},
	}
	tempDir := t.TempDir()
	for _, tCase := range testCase {
		if tCase.err != os.ErrNotExist {
			if _, err := os.Create(filepath.Join(tempDir, tCase.path)); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}
	}
	for _, tCase := range testCase {
		result, err := validPath(filepath.Join(tempDir, tCase.path))
		if result != tCase.result || err != tCase.err {
			t.Errorf("got result %v - err %v, wanted result %v - err %v", result, err, tCase.result, tCase.err)
		}
	}
}

type TestZipFolderCase struct {
	folderName    string
	files         []string
	zipFolderName string
}

func TestZipFolder(t *testing.T) {
	// arrange

	testCase := TestZipFolderCase{
		folderName: "testzip",
		files: []string{
			"a.go", "b.go", "c.go", "d.go",
		},
		zipFolderName: "test.zip",
	}
	tempDir := t.TempDir()
	if err := os.Mkdir(tempDir+"/"+testCase.folderName, 0o755); err != nil {
		t.Error("Create test zip folder fail")
	}
	for _, file := range testCase.files {
		filePath := filepath.Join(tempDir, testCase.folderName, file)
		if _, err := os.Create(filePath); err != nil {
			t.Errorf("Create file error %v", err)
		}
	}
	// act
	if err := zipFolder(tempDir+"/"+testCase.folderName, filepath.Join(tempDir, testCase.zipFolderName)); err != nil {
		t.Errorf("Zip folder err %v", err)
	}
	// assert
	if rs, _ := validPath(filepath.Join(tempDir, testCase.zipFolderName)); !rs {
		t.Error("Zip folder not found")
	}
}
