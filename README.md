# hetzner-machine-provider

This is a provider for gitlab custom runners. It uses the Hetzner Cloud API to create and delete servers for using them inside the ci.
So, the CI behaves like in github actions with their "per-job" VMs.

## Usage
Environment variables options used in ci config file:
- **HCLOUD_SERVER_TYPE**: The server type to use, for example `ccx12`, defaults to `auto`
- **HCLOUD_SERVER_ARCHITECTURE**: The architecture to use for the server, only being used if `HCLOUD_SERVER_TYPE` is set to `auto`, defaults to `amd64`
- **HCLOUD_SERVER_LOCATION**: The location to use, defaults to `fsn1`
- **HMP_SERVER_WAIT_DEADLINE**: The time to wait for the server to be ready, defaults to `5m`
- **HMP_ADDITIONAL_AUTHORIZED_KEYS**: Additional authorized keys to add to the server, defaults to `""`. Separate multiple keys with a newline (`\n`).

### Image Selection
You can set the image to use by setting the `image` property in the `.gitlab-ci.yml` file.
If you don't set it, it will default to `ubuntu-22.04`.

Also, some special image selectors are available:
  - `:latest`-Suffix: Will be used to filter the images and selects the one with the highest os version. Example: `ubuntu:latest`
  - `label#`-Prefix: Will be used to filter with label selectors. That is used for snapshots. The snapshot with the latest creation date will be selected. See [docs](https://docs.hetzner.cloud/#label-selector) for examples.

## Runner Configuration
You need to configure the following environment variable for your gitlab runner:
- **HCLOUD_TOKEN**: The API token for the Hetzner Cloud API, must have the permissions to create and delete servers

Furthermore, you need to configure the runner to use the custom executor. Here is an example configuration:
```toml
concurrent = 4
check_interval = 0
shutdown_timeout = 0

[session_server]
  session_timeout = 1800
  ...
[[runners]]
  ...
  executor = "custom"
  builds_dir = "/builds"
  cache_dir = "/cache"
  [runners.custom]
    config_exec = "/<path-to-hmp>/hmp"
    config_args = ["configure"]
    prepare_exec = "/<path-to-hmp>/hmp"
    prepare_args = ["prepare"]
    run_exec = "/<path-to-hmp>/hmp"
    run_args = ["exec"]
    cleanup_exec= "/<path-to-hmp>/hmp"
    cleanup_args = ["cleanup"]
  ...
```
