package file

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
)

type FileTarget struct {
	Path   string
	Format string
}

type FileFormat string

const (
	Env  FileFormat = "env"
	JSON FileFormat = "json"
)

func (f *FileTarget) Args() {
	flag.StringVar(&f.Path, "target-file-path", "", "Path to file where the secret contents would be saved and stored")
	flag.StringVar(&f.Format, "target-file-format", "json", "Format in which the contents of the new keystore would be written to the target file. Currently, supported format are 'json' and 'env'")
}

func (f FileTarget) Execute(oldKeyStore, newKeyStore target.KeyStore) error {
	var content string
	switch f.Format {
	case string(Env):
		content = renderKeyStoreToEnvFormat(newKeyStore)
	case string(JSON):
		var err error
		content, err = renderKeyStoreToJsonFormat(newKeyStore)
		if err != nil {
			return fmt.Errorf("failed to render the key store the JSON format: %w", err)
		}
	default:
		return fmt.Errorf("unknown target format '%s' found. Currently, allowed formats are 'env', 'json'", f.Format)
	}
	log.Debugf("Writing the following content to the file at %s: %s", f.Path, content)
	contentBytes := []byte(content)
	if err := os.WriteFile(f.Path, contentBytes, 0644); err != nil {
		return fmt.Errorf("error occurred while writing to the file at the path '%s': %w", f.Path, err)
	}
	return nil

}

func renderKeyStoreToEnvFormat(keyStore target.KeyStore) string {
	content := ""
	for field, value := range keyStore {
		content += fmt.Sprintf("%s=%s\n", field, value)
	}
	return strings.TrimSuffix(content, "\n")
}

func renderKeyStoreToJsonFormat(keyStore target.KeyStore) (string, error) {
	contentBytes, err := json.Marshal(keyStore)
	if err != nil {
		return "", fmt.Errorf("error occurred while marshalling the JSON out of the key store: %w", err)
	}
	return string(contentBytes), err
}
