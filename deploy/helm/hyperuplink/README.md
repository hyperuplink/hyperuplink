# Hyperuplink Helm chart

Deploys the [Hyperuplink](https://hyperup.link) forum/BBS server on Kubernetes.

> **This chart deploys the Hyperuplink binary only.** PostgreSQL, Redis/Valkey
> and (optionally) an S3-compatible object store are **not** included. They must
> already exist and be reachable from the cluster. Point the config at them
> before installing.

## Prerequisites

- Kubernetes >= 1.23
- Helm >= 3.8
- A reachable PostgreSQL database and a Redis/Valkey instance
- An Ingress controller with TLS (production mode forces secure-only cookies)

## Installing

```sh
# From a local checkout:
helm install my-hyperuplink ./deploy/helm/hyperuplink \
  --namespace hyperuplink --create-namespace \
  --set-file config=my-hyperuplink.toml
```

`--set-file config=...` loads your TOML from a file into the `config` value; you
can also edit `config` inline via a `-f values.yaml` overrides file.

## Configuration

Hyperuplink is configured through a TOML file. The chart renders the `config`
value into a Secret and starts the container with
`-c file://<configMountPath>/<configFileKey>`.

Because the file holds credentials (database password, S3 keys, SMTP/XMPP
passwords, OAuth secrets) it is stored in a `Secret`, never a `ConfigMap`.

### Bring your own Secret

To manage the config Secret out of band (e.g. via SealedSecrets / External
Secrets), set `existingConfigSecret` to its name and expose the file under the
key named by `configFileKey`:

```yaml
existingConfigSecret: hyperuplink-config
configFileKey: hyperuplink.toml
```

When rendered by the chart, a `checksum/config` pod annotation triggers a
rolling restart whenever the config changes. With `existingConfigSecret` you are
responsible for restarting pods after editing the Secret.

### Key values

| Key                             | Default                           | Description                                         |
| ------------------------------- | --------------------------------- | --------------------------------------------------- |
| `replicaCount`                  | `1`                               | Replicas. Keep at 1 with local `persistence` (RWO). |
| `image.repository`              | `ghcr.io/hyperuplink/hyperuplink` | Image repository.                                   |
| `image.tag`                     | `""`                              | Image tag; defaults to the chart `appVersion`.      |
| `config`                        | see `values.yaml`                 | Inline `hyperuplink.toml`.                          |
| `existingConfigSecret`          | `""`                              | Use an existing config Secret instead of `config`.  |
| `configFileKey`                 | `hyperuplink.toml`                | Key/filename of the config in the Secret.           |
| `configMountPath`               | `/etc/hyperuplink`                | Where the config Secret is mounted.                 |
| `containerPort`                 | `3000`                            | Must match `Web.Port` in the config.             |
| `strategy.type`                 | `Recreate`                        | Deployment strategy (RWO-safe default).             |
| `persistence.enabled`           | `true`                            | PVC backing the `Local` storage provider.           |
| `persistence.size`              | `10Gi`                            | PVC size.                                           |
| `service.type` / `service.port` | `ClusterIP` / `3000`              | Service exposure.                                   |
| `ingress.enabled`               | `false`                           | Ingress (enable TLS for production mode).           |
| `autoscaling.enabled`           | `false`                           | HPA. Requires S3 storage (not local).               |
| `podDisruptionBudget.enabled`   | `false`                           | PDB for the deployment.                             |
| `networkPolicy.enabled`         | `false`                           | NetworkPolicy (egress left open by default).        |
| `resources`                     | 1 CPU / 512Mi limit               | Container resources.                                |

See [`values.yaml`](./values.yaml) for the full, commented list.

## Storage

- **Local** (default): media is written to a `ReadWriteOnce` PVC mounted at
  `persistence.mountPath` (`/var/lib/hyperuplink/media`). It is per-pod and not
  shared, so keep `replicaCount: 1` and do not enable autoscaling.
- **S3**: add an `[[Storage]]` block of `Type = "S3"` to `config`, set
  `persistence.enabled: false`, and you may then scale horizontally
  (`replicaCount > 1`, `autoscaling`, `strategy.type: RollingUpdate`).

## Production mode & TLS

`Mode = "production"` serves the embedded views and forces secure-only cookies,
so the app must be reached over HTTPS. Terminate TLS at the Ingress (set
`ingress.tls`) or an upstream proxy, and set `Web.ProxyHeader` /
`Web.TrustProxy` in the config so client IPs are read from the proxy.

## Uninstalling

```sh
helm uninstall my-hyperuplink --namespace hyperuplink
```

PVCs created by the chart are retained by Kubernetes and must be deleted
manually if you want to discard the media data.

## Testing the release

```sh
helm test my-hyperuplink --namespace hyperuplink
```

Runs a short-lived pod that probes the readiness endpoint through the Service.
