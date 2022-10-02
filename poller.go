package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
)

func Poller(vaultClient *vault.Client, tg target.Target, interval time.Duration, path string, exit chan os.Signal, failureLimit int64) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	oldKeyStore := target.KeyStore{}
	remainingFailures := failureLimit
	for {
		select {
		case <-ticker.C:
			newKeyStore, err := renderKeyStore(vaultClient, path)
			if err != nil {
				err = fmt.Errorf("error occurred while getting the contents of the key store at the path '%s': %w", path, err)
				if remainingFailures == 0 {
					log.Warn("failure limit reached, exiting!")
					return err
				}
				remainingFailures--
				log.Warn(err.Error())
				break
			}
			if reflect.DeepEqual(oldKeyStore, newKeyStore) {
				remainingFailures = failureLimit
				break
			}
			log.Debug("change observed between the old and new keystore, proceeding to execute the target")
			err = tg.Execute(oldKeyStore, newKeyStore)
			if err != nil {
				err = fmt.Errorf("error occurred while executing the target: %w", err)
				if remainingFailures == 0 {
					log.Warn("failure limit reached, exiting!")
					return err
				}
				remainingFailures--
				log.Warn(err.Error())
			} else {
				remainingFailures = failureLimit
				oldKeyStore = newKeyStore
			}
		case sig := <-exit:
			log.Infoln("receiving the signal", sig, "Exiting!")
			return nil
		}
	}
}
