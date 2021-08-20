# Example

The following example will generate YAML and Markdown docs for
[Docker buildx](https://github.com/docker/buildx) CLI.

```console
git clone https://github.com/docker/cli-docs-tool
cd cli-docs-tool/example/
go mod download
go run main.go
```

Generated docs will be available in `./docs`
