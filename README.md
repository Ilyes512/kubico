# Kubico

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/ilyes512/kubico.svg)](https://hub.docker.com/r/ilyes512/kubico)
[![Docker Pulls](https://img.shields.io/docker/pulls/ilyes512/kubico.svg)](https://hub.docker.com/r/ilyes512/kubico)
[![MicroBadger Layers](https://img.shields.io/microbadger/layers/ilyes512/kubico.svg)](https://microbadger.com/images/ilyes512/kubico)

## How to use/build

Requirements:
- [Docker](https://docs.docker.com/install/)
- [Task](https://taskfile.dev/#/installation) (A Task runner)
- [Go](https://golang.org/doc/install)

```
# Run using golang on host
task run

# Install golang host deps needed for building on host
task dl-deps

# Build on host
task build

# For list of tasks
task --list
```

## Screenshot

<div align="center">
  <img width="500" src="docs/assets/images/kubico.png">
</div>
