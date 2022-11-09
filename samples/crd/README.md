# CRD Example

The following sample demonstrate how to use the ServiceMapper operator to generate a ProvisionedService (ServiceProxy) and ServiceEndpointDefinition for an instance of MongoDB.
MongoDB instances are provisioned through and backed by the [MongoDB community Operator](https://github.com/mongodb/mongodb-kubernetes-operator).

## Prerequisites

- minikube
- make
- helm
- kubectl

## Demo

### Start minikube
Use the script `hack/start-minikube.sh` to start minikube

```bash
./hack/start-minikube.sh
```

### Deploy the Service Mapper operator

Configure the current shell to use minikube's docker, build and deploy the operator.

```bash
eval $(minikube docker-env)
make docker-build
make deploy
```

### Install MongoDB community operator via Helm

Use Helm to install the MongoDB community operator into the cluster

```bash
helm repo add mongodb https://mongodb.github.io/helm-charts
helm install community-operator mongodb/community-operator --namespace srm-sample --create-namespace
```

### Apply demo resources

Apply the manifests for the demo.
A MongoDB instance will be published together with it's ServiceResourceMap.
Also, a ClusterRole and a ClusterRoleBinding will be provisioned to authorize the ServiceMapperOperator to read/list/watch MongoDB instances.

```bash
kubectl apply -f samples/crd/
```

### Inspect ServiceEndpointDefinitions and ServiceProxy

View ServiceProxy's details

```
kubectl get -n srm-sample serviceproxies.binding.operators.coreos.com mongodb-srm-sample -o yaml
```

View secret's details:

```
kubectl get -n srm-sample secrets mongodb-srm-sample-sed -o yaml
kubectl get -n srm-sample secrets mongodb-srm-sample-sed --output json | jq '.data | map_values(@base64d)'
```


