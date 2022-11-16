# RDS Demo with CRC

## Prerequisites

Tools:

- crc
- docker
- kubectl
- make
- oc

## Start CRC cluster

```
crc start
eval $(crc oc-env)
oc login -u kubeadmin https://api.crc.testing:6443
oc registry login --skip-check
```

## Install ACK Operator

Manifests for installing the ACK Operator are stored in the `operators/ack` folder.

Use the `ack-secret.yaml.tmpl` to create a `ack-scret.yaml` file with plaintext AWS Access Key Id and Access Key Secret.
Refer to [https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey) for creating Access Keys from AWS Console.

```
kubectl apply -f samples/rds/operators/ack
```

Install the "AWS Controllers for Kubernetes - Amazon RDS" from the "Operator Hub" (cfr. [https://developers.redhat.com/articles/2022/09/21/bind-services-created-aws-controllers-kubernetes#step_4___create_a_database_instance](https://developers.redhat.com/articles/2022/09/21/bind-services-created-aws-controllers-kubernetes#step_4___create_a_database_instance)).

To open the console you can use the following command:

```
crc console
```

1. If **Developer** is selected in the dropdown on the left, switch to **Administrator** view
2. Click on the **Operators** blade
3. Select **OperatorHub**
4. Use the search bar to find **AWS Controllers for Kubernetes - Amazon RDS**
5. Install the operator with default parameters

## Install Service Binding Operator

Install the Service Binding Operator from the "Operator Hub", as you have done for [ACK Operator](#install-ack-operator)

Register the ClusterRole to allow the Service Binding Operator to get,list,watch Service Proxies.

```
kubectl apply -f samples/rds/operators/sbo
```

## Prepare project's namespace

```
kubectl apply -f samples/rds/ack-rds-namespace.yaml
```

## Deploy RDS resources

Create a Postgres DBInstance with AWS RDS by applying the manifests stored into the `postgres` folder.

```
kubectl apply -f samples/rds/postgres
```

> To access the database you may need to configure the inbound rules for the security group associated to the database.

## Run the Service Mapper Operator into CRC

Build the image and push it to the cluster

```
make docker-build IMG=$(oc registry info)/service-mapper-system/srm-controller:latest
kubectl create namespace service-mapper-system
docker push $(oc registry info)/service-mapper-system/srm-controller:latest
```

Bake manifests and deploy them to the cluster

```
make deploy IMG=image-registry.openshift-image-registry.svc:5000/service-mapper-system/srm-controller:latest
```

## Deploy ServiceResourceMap

When the RDS Postgresql database and the ServiceMapper Operator are running, deploy the ServiceResourceMap

```
kubectl apply -f samples/rds/ack-rds-psql-serviceresourcemap.yaml
```

## Check for ServiceProxy and SED

View ServiceProxy's details

```
kubectl get -n srm-rds-sample serviceproxies.binding.operators.coreos.com srm-rds-psql-sample -o yaml
```

View secret's details:

```
kubectl get -n srm-rds-sample secrets srm-rds-psql-sample-sed -o yaml
kubectl get -n srm-rds-sample secrets srm-rds-psql-sample-sed --output json | jq '.data | map_values(@base64d)'
```

## Deploy an application and bind to the ServiceProxy

Now you can deploy an application, and use the ServiceProxy to bind it to the Service.

```
kubectl apply -f samples/rds/app
```

And inspect the volumens attached to the created pod

```
kubectl get pods -l app=srm-rds-sample-app --output jsonpath='{.items[0].spec.volumes[0]}'
```

Finally we can inspect the content of the binding folder by exec-ing into the pod's container

```
kubectl exec -it $(kubectl get pods -l app=srm-rds-sample-app --output jsonpath='{.items[0].metadata.name}') -- sh
```

Into the shell the previous command created, we can inspect the projected data

```
ls -lh /bindings/srm-rds-psql-sample

# try read some projected secrets
cat /bindings/srm-rds-psql-sample/host
```

