# hcloud-machine-provider

This is a provider for gitlab custom runners. It uses the Hetzner Cloud API to create and delete servers for using them inside the ci.

## Usage
You need to configure the following environment variables:
- **HCLOUD_TOKEN**: The API token for the Hetzner Cloud API, must have the permissions to create and delete servers

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
    prepare_exec = "/<path-to-hmp>/hmp"
    prepare_args = ["prepare"]
    run_exec = "/<path-to-hmp>/hmp"
    run_args = ["exec"]
    cleanup_exec= "/<path-to-hmp>/hmp"
    cleanup_args = ["cleanup"]
  ...
```