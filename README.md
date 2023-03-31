## Fly.io Gossip Gloomers Distributed Systems Challenges - Grow-only Counter Challenge

This repo contains a Go implementation of the grow-only counter challenge for the [Fly.io Gossip Gloomers](https://fly.io/dist-sys/) series of distributed systems challenges.

## Requirements

- Go 1.20: you can install it using [gvm](https://github.com/moovweb/gvm) with `gvm install go1.20 && gvm use go1.20`

## Build

From the project's root directory:

```bash
go build .
```

## Test

To use the different Maelstrom test commands, you can refer to the Fly.io [instructions](https://fly.io/dist-sys/4/), or run:

```bash
# Make sure to replace `~/go/bin/maelstrom-counter`
# with the full path of the executable you built above
./maelstrom test -w g-counter \
  --bin ~/go/bin/maelstrom-counter \
  --node-count 3 \
  --rate 100 \
  --time-limit 20 \
  --nemesis partition
```
