# Optimus task plugin - NEO tracker

Example task to demonstrate **near earth orbit** asteroids. It gives a list of potential
hazardous object in space based on date using [NASA API](https://api.nasa.gov/).

Plugins are implementations of Task and Hook interfaces for supporting the execution
of arbitrary jobs in optimus.

- This requires NASA API to work, register as a secret with the name `optimus-task-neo`
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

### TODO

- Find neo object images and print them using image to ascii converter