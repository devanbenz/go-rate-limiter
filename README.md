# Rate-limiter

Diagram
```
                                       ┌───────┐
                                       │       │
                             ┌─────────► redis │
                             │         │       │
                             │         └┬──────┘
                             │          │
                             │          │      ┌───────────────┐
                             │          │      │               │
                       ┌─────┴─────┐    │      │               │
┌──────────┐           │           ◄────┘      │               │
│          │           │           │           │    server     │
│  client  ├───────────► middleware├───────────►               │
│          │           │           │           │               │
└──────────┘           │           │           │               │
                       └───────────┘           │               │
                                               └───────────────┘
```

Pretty basic rate-limit built in an hour. It currently uses the token bucket algorithm to cache request
counts from the remote address. See: https://en.wikipedia.org/wiki/Token_bucket
