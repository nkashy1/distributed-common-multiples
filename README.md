# distributed-common-multiples

A distributed system for the calculation of common multiples of positive integers.

To play with this, first install [go](https://golang.org/).

Then,
```
go run main.go <n_1> <n_2> ... <n_m>
```
where each `<n_j>` is a positive integer.

Example:
```
$ go run main.go 1289031 2 3 5
.
.
.
Acceptor with key 2: rejected proposal 1199
Acceptor with key 1289031: rejected proposal 1199
Acceptor with key 5: rejected proposal 1199
Acceptor with key 3: rejected proposal 1199
Acceptor with key 1289031: rejected proposal 1200
Acceptor with key 2: accepted proposal 1200
Acceptor with key 3: accepted proposal 1200
Acceptor with key 5: accepted proposal 1200
CONSENSUS: 1200
Proposer: Killing self
```

The reason you are seeing `1200` and not `30` like you might is that the responses are being sent asynchronously and we are neither
using buffered channels nor blocking on channel sends.
