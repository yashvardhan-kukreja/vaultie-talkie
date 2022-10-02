package file

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
)

type mockDir struct {
	path string
}

func (m mockDir) setup() error {
	if m.path == "" {
		return fmt.Errorf("directory path not found to be provided")
	}
	return os.MkdirAll(m.path, os.ModePerm)
}

func (m mockDir) teardown() error {
	if m.path != "" {
		return os.RemoveAll(m.path)
	}
	return nil
}

func TestFileTargetSuccessWithJsonFormat(t *testing.T) {
	m := mockDir{
		path: "./test-filetarget",
	}

	assert.NoDirExists(t, m.path)
	assert.NoError(t, m.setup())
	assert.DirExists(t, m.path)

	ft := FileTarget{
		Path:   fmt.Sprintf("%s/sample.json", m.path),
		Format: "json",
	}
	oldKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
	})
	newKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
		"a":   "b",
	})
	assert.NoError(t, ft.Execute(oldKeyStore, newKeyStore))

	expectedContentsBytes, err := json.Marshal(newKeyStore)
	assert.NoError(t, err)

	writtenFileContentsBytes, err := os.ReadFile(ft.Path)
	assert.NoError(t, err)

	assert.Equal(t, expectedContentsBytes, writtenFileContentsBytes)

	assert.NoError(t, m.teardown())
	assert.NoDirExists(t, m.path)
}

func TestFileTargetSuccessWithEnvFormat(t *testing.T) {
	m := mockDir{
		path: "./test-filetarget",
	}

	assert.NoDirExists(t, m.path)
	assert.NoError(t, m.setup())
	assert.DirExists(t, m.path)

	ft := FileTarget{
		Path:   fmt.Sprintf("%s/sample.env", m.path),
		Format: "env",
	}
	oldKeyStore := target.KeyStore(map[string]interface{}{
		"a": "b",
	})
	newKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
	})
	assert.NoError(t, ft.Execute(oldKeyStore, newKeyStore))

	expectedContent := ""
	for k, v := range newKeyStore {
		expectedContent += fmt.Sprintf("%s=%s\n", k, v)
	}
	expectedContent = strings.TrimSuffix(expectedContent, "\n")

	writtenFileContentsBytes, err := os.ReadFile(ft.Path)
	assert.NoError(t, err)

	assert.Equal(t, expectedContent, string(writtenFileContentsBytes))

	assert.NoError(t, m.teardown())
	assert.NoDirExists(t, m.path)
}

func TestFileTargetFailureWithWrongFormat(t *testing.T) {
	m := mockDir{
		path: "./test-filetarget",
	}

	assert.NoDirExists(t, m.path)
	assert.NoError(t, m.setup())
	assert.DirExists(t, m.path)

	ft := FileTarget{
		Path:   fmt.Sprintf("%s/sample.json", m.path),
		Format: "gibberish",
	}
	oldKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
	})
	newKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
		"a":   "b",
	})
	assert.Error(t, ft.Execute(oldKeyStore, newKeyStore), fmt.Sprintf("unknown target format '%s' found. Currently, allowed formats are 'env', 'json'", ft.Format))

	assert.NoError(t, m.teardown())
	assert.NoDirExists(t, m.path)
}
