package shared_test

import (
	"errors"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/OpenSlides/openslides-manage-service/pkg/shared"
)

func TestCreateFile(t *testing.T) {
	t.Run("running shared.CreateFile() in empty directory", func(t *testing.T) {
		testDir, err := os.MkdirTemp("", "openslides-manage-service-")
		if err != nil {
			t.Fatalf("generating temporary directory failed: %v", err)
		}
		defer os.RemoveAll(testDir)
		fileName := "test_file_eeGhu6du_3"
		content := "test_content_kohv2EoT_3"

		shared.CreateFile(testDir, false, fileName, []byte(content), false)

		testContentFile(t, testDir, fileName, content)
	})

	t.Run("running shared.CreateFile() in non existing directory", func(t *testing.T) {
		testDir := "non_existing_directory"
		fileName := "test_file_eeGhu6du_4"
		content := "test_content_kohv2EoT_4"
		hasErrMsg := "no such file or directory"

		err := shared.CreateFile(testDir, false, fileName, []byte(content), false)

		if !strings.Contains(err.Error(), hasErrMsg) {
			t.Fatalf("running shared.CreateFile() with invalid directory, got error message %q, expected %q", err.Error(), hasErrMsg)
		}
	})

	t.Run("running shared.CreateFile() on existing file with force true", func(t *testing.T) {
		testDir, err := os.MkdirTemp("", "openslides-manage-service-")
		if err != nil {
			t.Fatalf("generating temporary directory failed: %v", err)
		}
		defer os.RemoveAll(testDir)
		fileName := "test_file_eeGhu6du_5"
		content := "test_content_kohv2EoT_5"
		shared.CreateFile(testDir, false, fileName, []byte(content), false)
		content = "test_content_kohv2EoT_5b"

		shared.CreateFile(testDir, true, fileName, []byte(content), false)

		testContentFile(t, testDir, fileName, content)
	})

	t.Run("running shared.CreateFile() on existing file with force false", func(t *testing.T) {
		testDir, err := os.MkdirTemp("", "openslides-manage-service-")
		if err != nil {
			t.Fatalf("generating temporary directory failed: %v", err)
		}
		defer os.RemoveAll(testDir)
		fileName := "test_file_eeGhu6du_6"
		content := "test_content_kohv2EoT_6"
		shared.CreateFile(testDir, false, fileName, []byte(content), false)
		content2 := "test_content_kohv2EoT_6b"

		err2 := shared.CreateFile(testDir, false, fileName, []byte(content2), false)

		if err2 != nil {
			t.Fatalf("running shared.CreateFile() with invalid directory, got error message %q, expected nil error", err2.Error())
		}
		testContentFile(t, testDir, fileName, content)

	})
}

func testContentFile(t testing.TB, dir, name, expected string) {
	t.Helper()

	p := path.Join(dir, name)
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("file %q does not exist, expected existance", p)
	}
	content, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("error reading file %q: %v", p, err)
	}

	got := string(content)
	if got != expected {
		t.Fatalf("wrong content of file %q, got %q, expected %q", p, got, expected)
	}
}