# Vaultie-Talkie

A tool which observes the secrets present in your Hashicorp Vault store and trigger events defined by you whenever those secrets are updated.

---

## I'll cut to the chase

* You are saving secrets in Hashicorp Vault.
* You or your application servers or some other piece of code of yours makes use of those secrets to work.
* Updating the secrets is a hassle for you.
    * Everytime you update the secret, you have to tell your application that the secret changed by restarting it.
    * Or you manually end up triggering some job/script to propagate the fact that your secret changed.
    * Or at times, you just wanna tell/convey your team or some other team that someone/something changed secrets in your vault store over Slack.

* Hashicorp vault doesn't produce any events or have any webhook support with which you can plugin some business logic handling the events of secrets changes.

* So what to do? You use vaultie-talkie

## Introduction

Vaultie-talk is an event recorder and webhook simulator for changes happening in your Hashicorp Vault store.

You can configure vaultie-talkie to do the following:
* Trigger a webhook you own whenever a secret is changed in your vault store. Vaultie-talkie will send the details of both the old and new state of secrets (before and after the change respectively) as a part of the request sent to the webhook.
* Update a file, say a json file or .env file, tracking the vault secrets in, say, your application servers.
* Or just execute a command whenever the vault secret changes. Just before executing the command, vaultie-talkie would save the details of old and new secret in a temporary file so that the command which runs could assume to reach out to that temporary file, expect and get the old and new secret details, and accordingly take any action it intends to. Now, this command could be literally anything. It can be a command which mails you "Hello, your secret changed."

## Installation

```sh
curl -L https://raw.githubusercontent.com/yashvardhan-kukreja/vaultie-talkie/main/scripts/install.sh | bash
```

## Usage

vaultie-talkie expects a mandatory set of arguments no matter how you use it.

And then, depending on the "target" type you choose, it would expect some other extra arguments.

When I say "target", I mean the action you expect vaultie-talkie to take whenever it observes a secret change i.e. one of the following:
- Trigger a webhook on secret change?
- Update a file with the new secret contents?
- Execute a command on secret change?

So, here we go:

### Mandatory Arguments:

| argument            | value type                                              | default | explanation                                                                                                                                                                                                                                           |   |
|---------------------|---------------------------------------------------------|---------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---|
| `-vault-host`         | string                                                  | ""      | Host (without the protocol) at which your vault store is running. For example, "23.45.67.89"                                                                                                                                                          |   |
| `-vault-port`         | int                                                     | 8200    | Port at which your vault store is running. For example, 8200                                                                                                                                                                                          |   |
| `-vault-path`         | string                                                  | ""      | Vault path of the secret you want vaultie-talkie to watch and respond to. For example, "foo/bar"                                                                                                                                                      |   |
| `-vault-access-token` | string                                                  | ""      | Vault token which has at least read privileges to the above vault path                                                                                                                                                                                                        |   |
| `-target-type`        | string (allowed values: "webhook" / "file" / "command") | ""      | Type of action which vaultie-talkie would take when the secret contents at -vault-path change.                                               |   |
| `-failure-limit`      | int                                                     | 3       | Amount of failures/errors the poller should be allowed bear in a row. Once the this number is reached, vaultie-talkie would exit. Until then, it's going to just log the errors and retry. For setting no/infinite failure limits, feed the value -1. |   |
| `-polling-interval`   | time                                                    | 5s      | Rate at which the vault store gets polled for watching its contents                                                                                                                                                                                   |   |
| `-debug`              | bool                                                    | false   | Run vaultie-talkie in debug mode. Would log extra logs in the console where vaultie-talkie would be running.                                                                                                                                          |   |

### Arguments for "webhook" target-type

| argument              | value type                                              | default | explanation                                                                                                                                                                                                                                           |   |
|-----------------------|---------------------------------------------------------|---------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---|
| `-webhook-url`          | string                                                  | ""      | Webhook URL against which a POST request would be triggered whenever the secret contents under -vault-path changes.                                                                                                                                   |   |
| `-webhook-access-token` | string                                                  | ""      | Authorization Bearer token to be used by vaultie-talkie while making requests against the Webhook URL.                                                                                                                                                |   |

### Arguments for "file" target-type

| argument            | value type                                              | default | explanation                                                                                                                                                                                                                                           |   |
|---------------------|---------------------------------------------------------|---------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---|
| `-target-file-path`   | string                                                  | ""      | Path to the file where the secret contents would be saved by vaultie-talkie.                                                                                                                                                                          |   |
| `-target-file-format` | string (allowed values: "json" / "env")                 | "json"  | Format in which the contents of the new secret would be written to the target file. Currently, supported format are 'json' and 'env'                                                                                                                  |   |

### Arguments for "command" target-type

| argument                                | value type                                              | default                             | explanation                                                                                                                                                                                                                                                                                                                            |   |
|-----------------------------------------|---------------------------------------------------------|-------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---|
| `-target-command`                         | string                                                  | ""                                  | Command to execute by vaultie-talkie whenever the secret contents under -vault-path are changed.                                                                                                                                                                                                                                       |   |
| `-intermediate-file-for-changed-keystore` | string                                                  | "/tmp/vaultie-talkie/keystore.json" | At this file, the old and new secret contents will be written in the JSON format: {'old_key_store': <old key store JSON>, 'new_key_store': <new key store JSON>}. Whenever your target command executes, it can assume that the contents of the old and new secret would be present in this intermediate path and accordingly, use it. |   |

## Examples

In case you got bored reading the above instructions, here are some direct examples to aid you in understanding the usage of vaultie-talkie intuitively.

### Situation 1 (depicting the "webhook" target type)

I have a vault store running at `http://myvault.yash.com` and under it, I have a secret at the path `applications/ecommerce` with the following contents.

```json
{
    "secret1": "value1"
}
```
(The token "mysecrettoken" allows read access to this secret)

I want to receive a mail whenever the contents under the above secret are changed.
To accomplish this, I have written a webhook endpoint at `http://mywebhook.yash.com/webhook` which I would expect to get triggered whenever the above secret is changed with the details/info associated with the changed secret.

So, I'll start vaultie-talkie with the following command

```sh
vaultie-talkie \
-vault-host=myvault.yash.com -vault-path=applications/ecommerce -vault-access-token=mysecrettoken -polling-interval=3s \
\
-target-type=webhook \
\
-webhook-url=http://mywebhook.yash.com/webhook
```

And I start this. Now, what's gonna happen is the following:
- vaultie-talkie will look for Vault at `http://myvault.yash.com`
- It will query the secret present under `applications/ecommerce` every 3 seconds.
- Whenever it would encounter a change in the above secret, it would make the following request:

```
POST http://mywebhook.yash.com/webhook

Request Body:
```
```json
{
    "old_key_store": {
        "secret1": "value1"
    },
    "new_key_store": {
        "secret1": "newvalue1",
        "secret2": "value2"
    }
}
```

And as I own the webhook, I can program it underneath to parse such incoming request and take any action like sending me emails whenever it is triggered.

### Situation 2 (depicting "command" target type)


I have a vault store running at `http://myvault.yash.com` and under it, I have a secret at the path `applications/ecommerce` with the following contents.

```json
{
    "secret1": "value1"
}
```
(The token "mysecrettoken" allows read access to this secret)

Now, my use-case is interesting. I am running an application somewhere which relies on a `config.json` to run. The config.json is meant to resemble the above secret in vault identically.

Which means whenever the secret contents in the above vault path are changed, I want the `config.json` to be updated resembling the new secret contents and my application is restarted.

To accomplish the above task, I have written a python script on the same host as my application. This python script, when executed, would update the config.json assuming the latest secret contents present at `/tmp/latest-vault-contents` and then, restart the server. And the command `python /home/path/config-updater.py` would run the above python script.

So, this is how I'll use the vaultie-talkie to achieve my use-case.

```sh
vaultie-talkie \
-vault-host=myvault.yash.com -vault-path=applications/ecommerce -vault-access-token=mysecrettoken \
\
-target-type=command \
\
-target-command="python /home/path/config-updater.py" -intermediate-file-for-changed-keystore "/tmp/latest-vault-contents"
```

And I start this. Now, what's gonna happen is the following:
* vaultie-talkie will look for Vault at `http://myvault.yash.com`
* It will query the secret present under `applications/ecommerce` every 5 seconds (default).
* Whenever it would encounter a change in the above secret, it would do two things:
    * Update the intermediate file `/tmp/latest-vault-contents` with the following stuff:

    ```json
    {
        "old_key_store": {
            "secret1": "value1"
        },
        "new_key_store": {
            "secret1": "newvalue1",
            "secret2": "value2"
        }
    }
    ```
    * Execute the command `python /home/path/config-updater.py`

And I'll be writing the python script `/home/path/config-updater.py` like this


```py
import json
import os

new_config_json_content = {}

## read the new contents of the secret from the intermediate file
with open("/tmp/latest-vault-contents") as f:
    vault_change = json.load(f)
    new_config_json_content = vault_change["new_key_store"]

## update config.json with the new secret
with open('config.json', 'w') as f:
    json.dump(new_config_json_content, f)

## restart my application using config.json
os.system("restart-my-application")
```

## Contributing

This project has a great scope of contributing to it in terms adding new kinds of targets to it. The one which come to my head at the moment are mail targets or slack targets ensuring that whenever a secret is changed a mail or a slack messaged is broadcasted respectively.

For making any contributions, please raise an issue associated with it and feel free to work on its PR once that issue is acknowledged by the maintainers (currently, me :P)

While working on the PR, please ensure to install `pre-commit` and set it up with respect to the project locally by running

```
pre-commit install
```

Moreover, you can ensure the resilience of your changes beforehand by running the existing tests via

```
make test
```

## Roadmap

* [ ] Define cooler targets like slack, mail, etc.
* [ ] Improve the docs (obviously) :3
* [ ] Make this project compatible with watching other storage backends like Hashicorp Consul KV (long shot goal).
