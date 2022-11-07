# CRD Example

```bash
# start minikube
./hack/start-minikube.sh
eval $(minikube docker-env)

# configure shell and deploy operator
make docker-build
make deploy

# install the community MongoDb operator
helm repo add mongodb https://mongodb.github.io/helm-charts
helm install community-operator mongodb/community-operator --namespace srm-sample --create-namespace

# apply the resources
kubectl apply -f samples/crd/
```


