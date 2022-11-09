# Service Mapper

This project will demonstrate how to create provisioned services via service resource maps.

## Description

This project introduces the following CRDs:
* **ServiceResourceMap**: defines the rules needed to generate a ServiceEndpointDefinition for a GroupVersionResource. ServiceResourceMaps are Cluster scoped.
    ```yaml
    apiVersion: binding.operators.coreos.com/v1alpha1
    kind: ServiceResourceMap
    metadata:
      annotations:
        name: srm-sample-postgresql
    spec:
      service_kind_reference:
        api_group: rds.services.k8s.aws/v1alpha1
        kind: dbinstances
      service_map:
        host: path={.status.endpoint.address}
        password: path={.spec.masterUserPassword.name},objectType=Secret,sourceKey=password
        port: path={.status.endpoint.port}
        type: path={.spec.engine}
    ```
* **ServiceProxy**: Namespaced resource that implements the ServiceBinding's specification for Provisioned Service.
    ```yaml
    apiVersion: binding.operators.coreos.com/v1alpha1
    kind: ServiceProxy
    metadata:
      name: srm-rds-psql-sample
      namespace: srm-rds-sample
    spec:
      service_instance:
        name: srm-rds-psql-sample
        namespace: srm-rds-sample
      service_resource_map: srm-sample-postgresql
    status:
      binding:
        name: srm-rds-psql-sample-sed
   ```

### Users Experience

**Administrator** creates a ServiceResourceMap, the **operator** looks for instances of the services referenced in the ServiceResourceMap and creates a ServiceProxy for each instance.
ServiceProxies are created in the same project/namespace of the service instance.

The **operator** also monitors for events on instances referenced by the published ServiceResourceMap, and creates/updates/deletes related ServiceProxies and ServiceEndpointDefinitions.

When a **Developer** creates an instance of the service in it's project/namespace, the **operator** will then create the ServiceProxy and ServiceEndpointDefinition in the same project/namespace.

The **Developer** can now use the ServiceProxy with the ServiceBindingOperator to bind an application to the service.


## Samples

The following samples are available in the `samples` folder:
- [MongoDB and Minikube](./samples/crd/README.md)
- [Amazon RDS and crc (Openshift 4)](./samples/rds/README.md)


## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/service-mapper:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/service-mapper:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
