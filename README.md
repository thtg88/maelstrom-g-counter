## Fly.io Gossip Gloomers Distributed Systems Challenges - Grow-only Counter Challenge

This repo contains a Go implementation of the grow-only counter challenge for the [Fly.io Gossip Gloomers](https://fly.io/dist-sys/) series of distributed systems challenges.

## Requirements

### Go 1.20

You can install Go 1.20 using [gvm](https://github.com/moovweb/gvm) with:

```bash
gvm install go1.20
gvm use go1.20
```

### Maelstrom

Maelstrom is built in [Clojure](https://clojure.org/) so you'll need to install [OpenJDK](https://openjdk.org/).

It also provides some plotting and graphing utilities which rely on [Graphviz](https://graphviz.org/) & [gnuplot](http://www.gnuplot.info/).

If you're using Homebrew, you can install these with this command:

```bash
brew install openjdk graphviz gnuplot
```

You can find more details on the [Prerequisites](https://github.com/jepsen-io/maelstrom/blob/main/doc/01-getting-ready/index.md#prerequisites) section on the Maelstrom docs.

Next, you'll need to download Maelstrom itself.

These challenges have been tested against [Maelstrom 0.2.3](https://github.com/jepsen-io/maelstrom/releases/tag/v0.2.3).

Download the tarball & unpack it.

You can run the maelstrom binary from inside this directory.

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
