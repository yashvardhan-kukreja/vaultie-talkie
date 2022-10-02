package commandexecutor

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
	"os"
	"os/exec"
)

type CommandExecutorTarget struct {
	Command          string
	IntermediateFile string
}

func (c *CommandExecutorTarget) Args() {
	flag.StringVar(&c.Command, "target-command", "", "Command to execute whenever the key store change is observed.")
	flag.StringVar(&c.IntermediateFile, "intermediate-file-for-changed-keystore", "/tmp/vaultie-talkie/keystore.json", "At this file, the old and new keystore will be written in the JSON format: {'old_key_store': <old key store JSON>, 'new_key_store': <new key store JSON>}. Whenever your target command executes, it can assume that the contents of the old and new keystore would be present in this intermediate path and accordingly, use it.")
}

func (c *CommandExecutorTarget) Execute(oldKeyStore, newKeyStore target.KeyStore) error {
	intermediateFileContents := map[string]target.KeyStore{
		"old_key_store": oldKeyStore,
		"new_key_store": newKeyStore,
	}
	intermediateFileContentsBytes, err := json.Marshal(intermediateFileContents)
	if err != nil {
		return fmt.Errorf("error occurred while marshalling the intermediate file contents into JSON format: %w", err)
	}
	log.Debug("Writing the old and new key store to the intemediate file at the path", c.IntermediateFile)
	if err := os.WriteFile(c.IntermediateFile, intermediateFileContentsBytes, 0644); err != nil {
		return fmt.Errorf("error occurred while writing to the intermediate file at the path '%s': %w", c.IntermediateFile, err)
	}

	log.Debug("Executing the following command", c.Command)
	if _, err := exec.Command("sh", "-c", c.Command).Output(); err != nil {
		return fmt.Errorf("error occurred while executing the target command '%s': %w", c.Command, err)
	}

	return nil
}
