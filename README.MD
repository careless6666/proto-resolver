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


## Configuration description

We have several way to configure how resolve files

### Path mode
```json
{
  "type": "path",
  "source": "/Users/Documents/Api/my.proto",
  "relativePath": "project_path",
  "version": "v1"
}
```

 - Type - path - approach how receive proto files, you can write file or folder (copy full directory recursively)
 - Source - Path to directory/file where stored your proto files
 - Relative path - where we put content from source 
 - Version - tech paramatr used caching (optional)


### Url mode

You should write url where stored proto files (raw mode!)

 - Type - url - approach how receive proto files, you can write file only
 - Source - url to file, you can write only file in raw mode, without extra html code
 - Relative path - where we put content from source
 - Version - tech parameter used caching (optional) 

### Git mode

You should write url with https format to git repository

- Type - git - approach how receive proto files, under the hood tool clone or pull git repository
- Source - git url to repository, write url in http(s) format
- Relative path - you can configure here special file or folder in git repository for resolve, if you put empty string tool will clone full proto files from source repository and place it in vendor.pb with git internet address
- Commit Revision - write 12 digits of commit, it will resolve proto files for this commit
- Tag - tool resolve proto files for this git tag
- Branch - tool resolve files from special branch

Git authorization:
 - ssh key mode - tool use native git client installed in your system, if you configured access to gitlab or github with ssh key, it will be used it
 - github api token - you can pass github_token parameter for clone with github api token
 - gitlab api token - you can pass gitlab_token parameter for clone with gitlab api token (you can also pass gitlab_domain paramater if you gitlab domain not contains word gitlab)

Pull modes:
 - pull enabled - pull for git resource, actual for branch mode
 - pull disabled - clone or get from cache git repo info









