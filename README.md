# Optimus Plugins

Plugins are implementations of Task and Hook interfaces for supporting the execution
of arbitrary jobs in [optimus](https://github.com/odpf/optimus).

## NEO tracker

Example task to demonstrate **near earth orbit** asteroids. It gives a list of potential
hazardous object in space based on date using [NASA API](https://api.nasa.gov/).


- This requires NASA API token to work, register as a secret with the name `optimus-task-neo`
as follows
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: optimus-task-neo
type: Opaque
data:
  key.json: base64encodedAPIkeygoeshere
```

## Installation

### To install this plugin on mac
```shell
brew install kushsharma/taps/optimus-plugins-kush
```

### To install this plugin on optimus server
```shell
curl -sL ${PLUGIN_RELEASE_URL} | tar xvz
chmod +x optimus-*
mv optimus-* /usr/bin/
```
where `${PLUGIN_RELEASE_URL}` would point to release package.
For e.g.
```shell
curl -sL https://github.com/kushsharma/optimus-plugins/releases/download/v0.0.3/optimus-plugins_0.0.3_linux_amd64.tar.gz | tar xvz
```

### TODO

- Find neo object images and print them using image to ascii converter