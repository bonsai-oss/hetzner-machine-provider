# hetzner-machine-provider

This is a provider for gitlab custom runners. It uses the Hetzner Cloud API to create and delete servers for using them inside the ci.
So, the CI behaves like in github actions with their "per-job" VMs.

## Usage
You need to configure the following environment variable for your gitlab runner:
- **HCLOUD_TOKEN**: The API token for the Hetzner Cloud API, must have the permissions to create and delete servers

Optional environment variables used in ci config:
- **HCLOUD_SERVER_TYPE**: The server type to use, defaults to `ccx21`
- **HCLOUD_SERVER_LOCATION**: The location to use, defaults to `fsn1`

You can set the image to use by setting the `image` property in the `.gitlab-ci.yml` file. If you don't set it, it will default to `ubuntu-22.04`.

You need to configure the following in the runner config:
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