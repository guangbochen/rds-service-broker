# User Provided Service Broker

User Provided Service Broker is an example
[Open Service Broker](https://www.openservicebrokerapi.org/)
for use demonstrating the Kubernetes
Service Catalog.

For more information,
[visit the Service Catalog project on github](https://github.com/kubernetes-incubator/service-catalog).

## Installing the Chart

To install the chart with the release name `rds-broker`:

```bash
$ helm install charts/rds-broker --name rds-broker --namespace rds-broker
```

## Uninstalling the Chart

To uninstall/delete the `rds-broker` deployment:

```bash
$ helm delete rds-broker
```

The command removes all the Kubernetes components associated with the chart and
deletes the release.

## Configuration

The following tables lists the configurable parameters of the User Provided
Service Broker

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image` | Image to use | `guangbo/rds-broker:v0.0.1` |
| `imagePullPolicy` | `imagePullPolicy` for the rds-broker | `Always` |

Specify each parameter using the `--set key=value[,key=value]` argument to
`helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be
provided while installing the chart. For example:

```bash
$ helm install charts/rds-broker --name rds-broker --namespace rds-broker \
  --values values.yaml
```
