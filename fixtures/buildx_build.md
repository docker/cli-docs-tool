# docker buildx build

<!---MARKER_GEN_START-->
Start a build

### Aliases

`build`, `b`

### Options

| Name | Description |
| --- | --- |
| [`--add-host stringSlice`](https://docs.docker.com/engine/reference/commandline/build/#add-entries-to-container-hosts-file---add-host) | Add a custom host-to-IP mapping (host:ip) |
| `--allow stringSlice` | Allow extra privileged entitlement, e.g. network.host, security.insecure |
| [`--build-arg stringArray`](https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg) | Set build-time variables |
| `--cache-from stringArray` | External cache sources (eg. user/app:cache, type=local,src=path/to/dir) |
| `--cache-to stringArray` | Cache export destinations (eg. user/app:cache, type=local,dest=path/to/dir) |
| `--compress` | Compress the build context using gzip |
| [`-f`](https://docs.docker.com/engine/reference/commandline/build/#specify-a-dockerfile--f), [`--file string`](https://docs.docker.com/engine/reference/commandline/build/#specify-a-dockerfile--f) | Name of the Dockerfile (Default is 'PATH/Dockerfile') |
| `--iidfile string` | Write the image ID to the file |
| `--label stringArray` | Set metadata for an image |
| `--load` | Shorthand for --output=type=docker |
| `--network string` | Set the networking mode for the RUN instructions during build |
| `-o`, `--output stringArray` | Output destination (format: type=local,dest=path) |
| `--platform stringArray` | Set target platform for build |
| `--push` | Shorthand for --output=type=registry |
| `--secret stringArray` | Secret file to expose to the build: id=mysecret,src=/local/secret |
| `--ssh stringArray` | SSH agent socket or keys to expose to the build (format: `default\|<id>[=<socket>\|<key>[,<key>]]`) |
| [`-t`](https://docs.docker.com/engine/reference/commandline/build/#tag-an-image--t), [`--tag stringArray`](https://docs.docker.com/engine/reference/commandline/build/#tag-an-image--t) | Name and optionally a tag in the 'name:tag' format |
| [`--target string`](https://docs.docker.com/engine/reference/commandline/build/#specifying-target-build-stage---target) | Set the target build stage to build. |


<!---MARKER_GEN_END-->

