package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
	commandExecutorTarget "github.com/yashvardhan-kukreja/vaultie-talkie/internal/target/commandexecutor"
	fileTarget "github.com/yashvardhan-kukreja/vaultie-talkie/internal/target/file"
	webhookTarget "github.com/yashvardhan-kukreja/vaultie-talkie/internal/target/webhook"
)

type Opts struct {
	VaultSettings
	TargetType      string
	PollingInterval time.Duration
	FailureLimit    int64
	DebugMode       bool
}

var validTargets = map[target.TargetType]target.Target{
	target.Webhook:         &webhookTarget.WebhookTarget{},
	target.File:            &fileTarget.FileTarget{},
	target.CommandExecutor: &commandExecutorTarget.CommandExecutorTarget{},
}

func main() {
	opts := Opts{}
	flag.StringVar(&opts.Host, "vault-host", "", "Host of the vault store backing your secrets")
	flag.Int64Var(&opts.Port, "vault-port", 8200, "Port at which the vault store is running")
	flag.StringVar(&opts.PathToWatch, "vault-path", "", "Path of the secret in the vault store to watch")
	flag.StringVar(&opts.AccessToken, "vault-access-token", "", "Access token authorizing to read/list the above path")
	flag.StringVar(&opts.TargetType, "target-type", "", "Type of action to happen upon vault key changes")
	flag.Int64Var(&opts.FailureLimit, "failure-limit", 3, "Amount of failures/errors the poller should bear in a row. Once the this number is reached, vaultie-talkie would exit. Until then, it's going to just log the errors and retry. For setting no/infinite failure limits, feed the value -1.")
	flag.DurationVar(&opts.PollingInterval, "polling-interval", 5*time.Second, "Rate at which the vault store gets polled for watching its contents")
	flag.BoolVar(&opts.DebugMode, "debug", false, "Run vaultie-talkie in debug mode")
	for _, tg := range validTargets {
		tg.Args()
	}
	flag.Parse()

	if opts.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	tg, ok := validTargets[target.TargetType(opts.TargetType)]
	if !ok {
		log.Fatal("unknown target type found")
	}

	log.Debug("parsed options", opts)
	log.Debug("parsed target's options", tg)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGHUP)

	vaultClient, err := opts.InitClient()
	if err != nil {
		log.Fatal("failed to initialize the vault client as per the provided parameters: %w", err)
	}

	log.Debug("vault client setup successfully")
	log.Debug("starting the poller...")

	if err := Poller(vaultClient, tg, opts.PollingInterval, opts.PathToWatch, exit, opts.FailureLimit); err != nil {
		log.Fatal(err)
	}
}
