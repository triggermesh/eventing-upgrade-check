# Eventing Upgrade Check v0.13.x --> v0.14.x

There are two functions that generate names for objects at Knative base code, one of them is being deprecated and the code being migrated to use the other.

Generated names wont match between parent objects names are long, causing an object re-creation along the way. These objects might be affected when upgrading eventing controller to v0.14.x:

- APIServer Source adapter deployment
- Ping Source adapter deployment
- Subscriptions

When recreating the items above you might experiment a temporary downtime, but for `Subscriptions` state might also be lost. `Subscriptions` are binded to channels, creating a client at the underlying software which might track offsets. That information will be delete when the current subscription is deleted.

Refer to the underlying channel software on how to flush operate to avoid losing or receiving duplicate events.

- Issue: https://github.com/knative/eventing/issues/2842
- PR:  https://github.com/knative/eventing/pull/2861

## What is this

This tool returns a text report with all subscriptions that will be re-created upgrading Knative Eventing from v0.13.x to v0.14.x.

It iterates triggers looking for their subscriptions and once found it will generate previous and new names, if they don't match they will be added to the output report.

This tool does not modify any object at your cluster.

## Usage

### Binary

`eventing-upgrade-check` looks for two environment variables on start up:

- `KUBECONFIG`, points to the location of the kubeconfig to use.
- `NAMESPACE`, is the namespace that will be checked for Triggers.

If `KUBECONFIG` is not found it will try to use in cluster configuration.
If `NAMESPACE` is not found it will check all Triggers at all namespaces.

- Download the binary for your OS from the [releases](https://github.com/triggermesh/triggerflow/releases) links.
- Set execution permissions
- Set environment variables if needed

Example execution

```sh
$ NAMESPACE=odacremolbap KUBECONFIG=/home/pablo/.kube/config  ./eventing-upgrade-check-linux-amd64
Starting upgrade-check v0.13.x to v0.14.x.
Found 1 subscriptions that need upgrade.
Subscription needs update:
        namespace: odacremolbap
        old name: text-message-hub-sendme-ev-b67e7354-3bc6-4cd6-b45d-201d8f477e5d
        new name: text-message-hub-sendme-everyth7125a19ae9c18d2a713dcdb5b84ccadb
        found: true
```

The report above indicates that one subscription that needs upgrading was found, printing the current new and the one that will be generated after v0.14.x eventing controller is started.

Old and new names are generated using Triggers and Brokers data, this tool also looks for the existing subscription and reports `found: true` if it has been retrieved. If `found: false` it might mean that it doesn't exist or there was an error retrieving the subscription in which case a log should have been written.

### Image

An image for this tool is provided at `gcr.io/triggermesh/eventing-upgrade-check:v0.1.0`

A set of manifest to run the update check inside the cluster is also provided. A job containing the binary is executed at the `knative-upgrade` namespace under the `upgrade-check` account with minimal read permissions on triggers and subscriptions.

To deploy the upgrade check using provided manifest run:

```sh
kubectl apply -f https://raw.githubusercontent.com/triggermesh/eventing-upgrade-check/master/deploy/all-in-one.yaml
```

You can then check the report:

```sh
$ kubectl logs -n knative-upgrade -l app=upgrade-check
Subscription needs update:
        namespace: metr1ckzu
        old name: default-test-trigger-from--56d31d3e-32f9-11ea-be82-42010a800192
        new name: default-test-trigger-from-yaml-653c61f8de5d23c94739755596ff8e6a
        found: true
Subscription needs update:
        namespace: zelig880
        old name: text-message-hub-sendme-ev-b67e7354-3bc6-4cd6-b45d-201d8f477e5d
        new name: text-message-hub-sendme-everyth7125a19ae9c18d2a713dcdb5b84ccadb
        found: true
```

Once you are done you can remove all items deployed:

```sh
kubectl delete -f https://raw.githubusercontent.com/triggermesh/eventing-upgrade-check/master/deploy/all-in-one.yaml
```

## Build

A Makefile is provided with a simplified set of targets

```sh
# Compile the binary
make build

# Compile binary for all OSes
make release

# Create container image (check variables to override default repo)
make image

# Push container image (check variables to override default repo)
make push

```

If using `ko` a config directory can be used to generate and run manifests.

```sh
ko create -f config
```
