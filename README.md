# Patterns

Common patters in go that I don't want to have to re-write for every project.

## Broadcasts

Often times, you might want everything that goes in 1 channel to be broadcast
to a bunch of other channels, optionally closing the output channels when
the input channel is closed. For an example which demonstrates this for strings,
see [broadcast-strings](./examples/broadcast-strings/main.go)
