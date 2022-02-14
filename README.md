# simple-cloud-provider

[simple-cloud-provider](https://github.com/dzeromski/simple-cloud-provider) is a general purpose cloud provider ;-)

## TL;DR

```console
# Add simple-cloud-provider repository
helm repo add simple-cloud-provider https://dzeromsk.github.io/simple-cloud-provider/
# Deploy a Helm Release named "my-release" using the simple-cloud-provider chart
helm install simple-cloud-provider simple-cloud-provider/simple-cloud-provider
```

## Introduction

This chart bootstraps a [simple-cloud-provider](https://github.com/dzeromski/simple-cloud-provider) deployment on a [Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Installing the Chart

To install the [Chart](https://helm.sh/docs/intro/using_helm/#three-big-concepts) with the [Release](https://helm.sh/docs/intro/using_helm/#three-big-concepts) name `simple-cloud-provider`:

```console
helm repo add simple-cloud-provider https://dzeromsk.github.io/simple-cloud-provider/
helm install simple-cloud-provider simple-cloud-provider/simple-cloud-provider
```

The command deploys simple-cloud-provider on the Kubernetes cluster in the default configuration.
The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `simple-cloud-provider` deployment:

```console
helm delete simple-cloud-provider
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

Please refer to [values.yaml](https://github.com/dzeromski/simple-cloud-provider/blob/master/aio/deploy/helm-chart/simple-cloud-provider/values.yaml)
for valid values and their defaults.

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```console
helm install simple-cloud-provider/simple-cloud-provider --name simple-cloud-provider \
  --set=service.externalPort=8080,resources.limits.cpu=200m
```

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
helm install simple-cloud-provider/simple-cloud-provider --name simple-cloud-provider -f values.yaml
```

> **Tip**: You can use the default [values.yaml](values.yaml), which is used by default, as reference

## Global and namespace pools

### Global pool

Any service in any namespace will take an address from the global pool `cidr/range`-global. 

### Namespace pool

A service will take an address based upon its namespace pool `cidr/range`-`namespace`. These would look like the following:

```
$ kubectl get configmap -n kube-system simple-cloud-provider -o yaml

apiVersion: v1
kind: ConfigMap
metadata:
  name: simple-cloud-provider
  namespace: kube-system
data:
  cidr-default: 192.168.0.200/29
  cidr-development: 192.168.0.210/29
  cidr-finance: 192.168.0.220/29
  cidr-testing: 192.168.0.230/29
```

## Create an IP pool using a CIDR

```
kubectl create configmap --namespace kube-system simple-cloud-provider --from-literal cidr-global=192.168.0.220/29
```

## Create an IP range

```
kubectl create configmap --namespace kube-system simple-cloud-provider --from-literal range-global=192.168.0.200-192.168.0.202
```

## Multiple pools or ranges

We can apply multiple pools or ranges by seperating them with commas.. i.e. `192.168.0.200/30,192.168.0.200/29` or `192.168.0.10-192.168.0.11,192.168.0.10-192.168.0.13`
