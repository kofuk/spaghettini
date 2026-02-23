# Spaghettini

A lightweight HTTP server that uses Go templates to generate responses.

Designed for cases where you need a little more than a static response, but far less than a full web framework.
Think stub servers, mock APIs, simple Webhook receivers, or any situation where you want to expose a quick HTTP interface with minimal logic baked in.

## Template Variables

The following variables are available inside templates:

| Variable | Type | Description |
|---|---|---|
| `.Request.Method` | `string` | HTTP method (`GET`, `POST`, ...) |
| `.Request.Path` | `string` | Request path |
| `.Request.Header` | `http.Header` | Request headers |
| `.Request.Body` | `string` | Raw request body |

Standard Go `text/template` syntax is supported, including `{{ if }}`, `{{ range }}`, `{{ with }}`, and so on.

## Examples

```shell
$ spaghettini --source examples/simple.gotpl
```

## License

MIT
