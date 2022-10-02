package target

type TargetType string

const (
	Webhook         TargetType = "webhook"
	File            TargetType = "file"
	CommandExecutor TargetType = "command"
)

type Target interface {
	Args()
	Execute(oldKeyStore, newKeyStore KeyStore) error
}

type KeyStore map[string]interface{}
