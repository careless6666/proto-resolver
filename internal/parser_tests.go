package main

//mock data

/* happy path
version: v1
deps:
  - git: github.com/googleapis/googleapis/google/api/http.proto v0.0.0-20211005231101-409e134ffaac
  - git: github.com/googleapis/googleapis/google/api/annotations.proto v0.0.0-20211005231101-409e134ffaac
  - remote_file: https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto
  - local_file: /tmp/path
*/

/*
variations of format

deps:
  - github.com/googleapis/googleapis/google/api v0.0.0-20211005231101-409e134ffaac
  - github.com/googleapis/googleapis/google/api/* v0.0.0-20211005231101-409e134ffaac
  - github.com/googleapis/googleapis/google/api/* v1.0.0

*/
