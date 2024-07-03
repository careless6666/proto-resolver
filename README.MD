# ProtoResolver

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


