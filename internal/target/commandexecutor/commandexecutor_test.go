package commandexecutor

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
)

// no-op

type mockFile struct {
	dirPath  string
	fileName string
}

func (m mockFile) setup() error {
	if m.dirPath == "" || m.fileName == "" {
		return fmt.Errorf("directory path or/and file name not found to be provided")
	}
	err := os.MkdirAll(m.dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.Create(path.Join(m.dirPath, m.fileName))
	if err != nil {
		m.teardown()
		return err
	}
	return nil
}

func (m mockFile) teardown() error {
	if m.dirPath != "" {
		return os.RemoveAll(m.dirPath)
	}
	return nil
}

func TestCommandExecutorSuccess(t *testing.T) {
	intermediateFile := mockFile{
		dirPath:  "./test-commandexecutor",
		fileName: "intermediate-file.json",
	}

	intermediateFilePath := path.Join(intermediateFile.dirPath, intermediateFile.fileName)
	assert.NoFileExists(t, intermediateFilePath)
	assert.NoError(t, intermediateFile.setup())
	assert.FileExists(t, intermediateFilePath)

	commandOutputFile := path.Join("./test-commandexecutor", "cmd-output.txt")

	ct := CommandExecutorTarget{
		Command:          "echo 'hello, world!' >> " + commandOutputFile,
		IntermediateFile: intermediateFilePath,
	}

	oldKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
	})
	newKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
		"a":   "b",
	})
	assert.NoError(t, ct.Execute(oldKeyStore, newKeyStore))

	intermediateFileContentsBytes, err := os.ReadFile(intermediateFilePath)
	assert.NoError(t, err)

	expectedIntermeditateFileContentsBytes, err := json.Marshal(map[string]interface{}{
		"old_key_store": oldKeyStore,
		"new_key_store": newKeyStore,
	})

	assert.NoError(t, err)
	assert.Equal(t, expectedIntermeditateFileContentsBytes, intermediateFileContentsBytes)

	commandOutputFileContentsBytes, err := os.ReadFile(commandOutputFile)
	assert.NoError(t, err)

	expectedCommandOutputFileContentsBytes := []byte("hello, world!\n")

	assert.Equal(t, string(expectedCommandOutputFileContentsBytes), string(commandOutputFileContentsBytes))

	assert.NoError(t, intermediateFile.teardown())
	assert.NoFileExists(t, intermediateFilePath)
}

func TestCommandExecutorFailureWithWrongCommand(t *testing.T) {
	intermediateFile := mockFile{
		dirPath:  "./test-commandexecutor",
		fileName: "intermediate-file.json",
	}

	intermediateFilePath := path.Join(intermediateFile.dirPath, intermediateFile.fileName)
	assert.NoFileExists(t, intermediateFilePath)
	assert.NoError(t, intermediateFile.setup())
	assert.FileExists(t, intermediateFilePath)

	ct := CommandExecutorTarget{
		Command:          "unknown-command-yada-yada",
		IntermediateFile: intermediateFilePath,
	}

	oldKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
	})
	newKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
		"a":   "b",
	})
	assert.Error(t, ct.Execute(oldKeyStore, newKeyStore), fmt.Sprintf("error occurred while executing the target command '%s': exit status 127", ct.Command))

	assert.NoError(t, intermediateFile.teardown())
	assert.NoFileExists(t, intermediateFilePath)
}
