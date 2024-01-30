# About

[![Build status](https://github.com/rgl/spin-http-go-example/workflows/build/badge.svg)](https://github.com/rgl/spin-http-go-example/actions?query=workflow%3Abuild)

Example Spin HTTP Application written in Go.

# Usage

Install [Go](https://github.com/golang/go), [TinyGo](https://github.com/tinygo-org/tinygo) and [Spin](https://github.com/fermyon/spin).

Start the application:

```bash
spin up --build
```

Access the HTTP endpoint:

```bash
xdg-open http://localhost:3000
```

# Container Image Usage

There are two ways to run a WebAssembly binary in a container:

1. Ship the `spin` binary and your `.wasm` file in a docker container image.
2. Ship the `.wasm` file and the `spin.toml` manifest file in a oci image
   artifact; this can be done in a single step with `spin registry push`.
   Then use a container runtime or orchestrator that supports running wasm
   containers. For example, `containerd` and the `containerd-shim-spin-v2`
   containerd shim.

**NB** The Fermyon Cloud directly uses the `.wasm` file, so there is no need to
use a container. Instead use `spin deploy` to deploy the application to the
Fermyon Cloud.

## Kubernetes Usage

See https://developer.fermyon.com/spin/v2/kubernetes.

## containerd Usage

Install [containerd](https://github.com/moby/containerd) and the [containerd-shim-spin](https://github.com/deislabs/containerd-wasm-shims/tree/main/containerd-shim-spin).

### containerd crictl Usage

Use `crictl`:

```bash
# see https://kubernetes.io/docs/concepts/architecture/cri/
# see https://kubernetes.io/docs/tasks/debug/debug-cluster/crictl/
# see https://kubernetes.io/docs/reference/tools/map-crictl-dockercli/
# see https://github.com/kubernetes-sigs/cri-tools/blob/master/docs/crictl.md
# see https://github.com/kubernetes-sigs/cri-tools/blob/master/cmd/crictl/sandbox.go
# see https://github.com/kubernetes-sigs/cri-tools/blob/master/cmd/crictl/container.go
# see https://github.com/kubernetes/cri-api/blob/kubernetes-1.27.10/pkg/apis/runtime/v1/api.proto
crictl pull \
  ghcr.io/rgl/spin-http-go-example:0.1.0
crictl images list
crictl info | jq .config.containerd.runtimes
install -d -m 700 /var/log/cri
cat >cri-spin-http-go-example.pod.yml <<'EOF'
metadata:
  uid: cri-spin-http-go-example
  name: cri-spin-http-go-example
  namespace: default
log_directory: /var/log/cri/cri-spin-http-go-example
EOF
cat >cri-spin-http-go-example.web.ctr.yml <<'EOF'
metadata:
  name: web
image:
  image: ghcr.io/rgl/spin-http-go-example:0.1.0
command:
  - /
log_path: web.log
EOF
pod_id="$(crictl runp \
  --runtime spin \
  cri-spin-http-go-example.pod.yml)"
web_ctr_id="$(crictl create \
  $pod_id \
  cri-spin-http-go-example.web.ctr.yml \
  cri-spin-http-go-example.pod.yml)"
crictl start $web_ctr_id
web_ctr_ip="$(crictl inspectp $pod_id | jq -r .status.network.ip)"
wget -qO- "http://$web_ctr_ip"
crictl ps -a                    # list containers.
crictl inspect $web_ctr_id | jq # inspect container.
crictl logs $web_ctr_id         # dump container logs.
crictl pods                     # list pods.
crictl inspectp $pod_id | jq    # inspect pod.
crictl stopp $pod_id            # stop pod.
crictl rmp $pod_id              # remove pod.
rm -rf /var/log/cri/cri-spin-http-go-example
```

### containerd ctr Usage

**NB** This is not yet working. See https://github.com/deislabs/containerd-wasm-shims/issues/202.

Use `ctr`:

```bash
ctr image pull \
  ghcr.io/rgl/spin-http-go-example:0.1.0
ctr images list
ctr run \
  --detach \
  --runtime io.containerd.spin.v2 \
  --net-host \
  ghcr.io/rgl/spin-http-go-example:0.1.0 \
  ctr-spin-http-go-example
ctr sandboxes list # aka pods.
ctr containers list
ctr container rm ctr-spin-http-go-example
```

# References

* [Spin Go SDK](https://github.com/fermyon/spin/tree/main/sdk/go)
* [Spin Go SDK Examples](https://github.com/fermyon/spin-js-sdk/tree/main/examples)
* [Spin Go SDK package](https://pkg.go.dev/github.com/fermyon/spin/sdk/go/v2)
* [Building Spin Components in Go](https://developer.fermyon.com/spin/v2/go-components)
* [Done icon](https://icons8.com/icon/uw-X2j32n7Xp/done)
* [Creating a container image](https://github.com/deislabs/containerd-wasm-shims/blob/main/containerd-shim-spin/quickstart.md#creating-a-container-image)
* [TinyGo WASI support](https://tinygo.org/docs/guides/webassembly/wasi/)
* [Go WASI support](https://go.dev/blog/wasi)
  * **NB** Upstream Go is not yet supported by Spin.
