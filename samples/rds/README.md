# RDS Demo

## Prerequisites

- Openshift 4.11/4.12

## Install ACK Operator

All manifests for installing the ACK Operator are stored in the `operators/ack` folder.

cfr. [https://developers.redhat.com/articles/2022/09/21/bind-services-created-aws-controllers-kubernetes#step_2___install_the_rds_operator_in_an_openshift_cluster](https://developers.redhat.com/articles/2022/09/21/bind-services-created-aws-controllers-kubernetes#step_2___install_the_rds_operator_in_an_openshift_cluster)

Use the `ack-secret.yaml.tmpl` to create a `ack-scret.yaml` file with plaintext AWS Access Key Id and Access Key Secret.
Refer to [https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey) for creating Access Keys from AWS Console.


## Install Service Binding Operator

All manifests for installing the Service Binding Operator are stored in the `operators/sbo` folder.


## Prepare project's namespace

```
kubectl apply -f ack-rds-namespace.yaml
```

## Deploy RDS resources

Create a Postgres DBInstance with AWS RDS by applying the manifests stored into the `postgres` folder.

## Run the Service Mapper Operator into CRC

```
crc start
eval $(crc oc-env)
oc registry login --skip-check
make docker-build IMG=$(oc registry info)/sbo-1225-system/srm-controller:latest
docker push $(oc registry info)/sbo-1225-system/srm-controller:latest
make deploy IMG=image-registry.openshift-image-registry.svc:5000/sbo-1225-system/srm-controller:latest
```


## Deploy ServiceResourceMap

When the RDS Postgresql database and the ServiceMapper Operator are running, deploy the ServiceResourceMap

```
kubectl apply -f ack-rds-psql-serviceresourcemap.yaml
```

## Check for ServiceProxy and SED

```
kubectl get -n srm-rds-sample serviceproxies.binding.operators.coreos.com
kubectl get -n srm-rds-sample secrets srm-rds-psql-sample-sed
```

View secret's data

```
kubectl get -n srm-rds-sample secrets srm-rds-psql-sample-sed --output json | jq '.data | map_values(@base64d)'
```
