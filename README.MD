# proto-resolver

## How to use it

Firstly you should configure file with name proto_deps.json with deps in format: 

http mode

```json
{
  "version": "v1",
  "deps": [
    {
      "type": "url",
      "source": "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto",
      "relativePath": "github.com/googleapis/googleapis/google/api",
      "version": "v1"
    }
  ]
}
```

Secondly install tool:

```shell
go install github.com/careless6666/proto-resolver/cmd/proto-resolver@latest
```

If that fails, make sure your GOPATH/bin is in your PATH. You can add it with:

```shell
export PATH=$PATH:$(go env GOPATH)/bin
```

In path with proto_deps.json execute 
```shell
proto-resolver restore
```
when you will see additional actions in console or result in vendor.pb directory




