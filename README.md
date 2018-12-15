# Kubia

This is a Dockerfile that builds a "Kubia"-image, the example image used in the book "Kubernetes in Action", except instead of using NodeJS it uses GoLang.

# Instructions

```bash
$ ./kubia -help
Usage of ./kubia:
  -unhealthyAfter int
        set the number of request after which the service should fail (returning 500 httpcode

$ ./kubia
Kubia server starting on port 'http://localhost:8080'...

# or
$ ./kubia -unhealthAfter 5
Kubia server starting in unhealthy mode on port 'http://localhost:8080'...
```
