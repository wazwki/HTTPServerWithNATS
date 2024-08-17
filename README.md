# How to config NATS in golang HTTP-server

## Install and import NATS library:

```bash
go get github.com/nats-io/nats.go
```
```go
import (
    "github.com/nats-io/nats.go"
)
```

## Basic Concepts

### Connection:

#### Connecting to the server: A `nats.Conn` object is created, representing a connection to the NATS server.

```go
nc, err := nats.Connect("nats://localhost:4222")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```

### Subscription:

#### Standard subscription: Allows subscribing to specific subjects and receiving messages.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Received message: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}
```

#### Channel subscription: Messages are placed in a channel, simplifying their processing in multithreaded applications.

```go
ch := make(chan *nats.Msg, 64)
sub, err := nc.ChanSubscribe("subject", ch)
if err != nil {
    log.Fatal(err)
}

go func() {
    for msg := range ch {
        log.Printf("Received message through channel: %s", string(msg.Data))
    }
}()
```

#### Queue subscription: Used for load balancing, where multiple subscribers receive messages from the same queue.

```go
sub, err := nc.QueueSubscribe("subject", "worker-group", func(m *nats.Msg) {
    log.Printf("Worker received: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}
```

### Publishing:

#### Simple publishing: Sends a message to the specified subject.

```go
err := nc.Publish("subject", []byte("Hello, NATS!"))
if err != nil {
    log.Fatal(err)
}
```

#### Publishing with acknowledgment (Ack): Waits for a confirmation of message delivery from the server.

```go
nc.Publish("subject", []byte("Hello, NATS!"))
err = nc.Flush()
if err != nil {
    log.Fatal(err)
}
log.Println("Message delivered!")
```

### Request-Reply:

#### Request: Sends a message and waits for a reply within a specified time.

```go
msg, err := nc.Request("subject", []byte("Can I get a response?"), 2*time.Second)
if err != nil {
    log.Fatal(err)
}
log.Printf("Received reply: %s", string(msg.Data))
```

#### Reply: Handles requests and sends a response back to the client.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Received request: %s", string(m.Data))
    m.Respond([]byte("Sure, here is your response!"))
})
if err != nil {
    log.Fatal(err)
}
```

### Auto-Unsubscribe:

#### You can set a limit on the number of messages a subscription can receive, after which it is automatically unsubscribed.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Received: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}

sub.AutoUnsubscribe(10)
```

## Advanced Features

### Connecting to a Cluster:

#### NATS supports connecting to multiple servers, which increases fault tolerance.

```go
nc, err := nats.Connect("nats://server1:4222,nats://server2:4222")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```

### Connection and Error Handling:

#### `nats.Conn` provides handlers for various events, such as disconnection, reconnection, error detection, and connection closure.

```go
nc, err := nats.Connect("nats://localhost:4222",
    nats.DisconnectHandler(func(_ *nats.Conn) {
        log.Println("Disconnected!")
    }),
    nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Printf("Reconnected to %v!\n", nc.ConnectedUrl())
    }),
    nats.ClosedHandler(func(_ *nats.Conn) {
        log.Println("Connection closed!")
    }),
)
```

### Durable Subscription:

#### NATS Streaming (an extended version of NATS) supports durable subscriptions that maintain state between sessions.

```go
import stan "github.com/nats-io/stan.go"

sc, err := stan.Connect("test-cluster", "client-id")
if err != nil {
    log.Fatal(err)
}
defer sc.Close()

sub, err := sc.Subscribe("subject", func(msg *stan.Msg) {
    log.Printf("Received a message: %s\n", msg)
}, stan.DurableName("my-durable"))
if err != nil {
    log.Fatal(err)
}
```

### Streaming Subscriptions:

#### The ability to store messages for subscribers that were temporarily unavailable and deliver them the missed data upon reconnection.

```go
import stan "github.com/nats-io/stan.go"

sc, err := stan.Connect("test-cluster", "subscriber")
if err != nil {
    log.Fatal(err)
}
defer sc.Close()

_, err = sc.Subscribe("orders", func(msg *stan.Msg) {
    log.Printf("Received: %s", msg.Data)
}, stan.StartWithLastReceived())
if err != nil {
    log.Fatal(err)
}

select {}
```

### Authorization and Encryption:

#### Support for authentication using tokens, passwords, and certificates.
#### TLS connections to ensure the security of data transmission.

```go
nc, err := nats.Connect("nats://localhost:4222", nats.Secure(&tls.Config{...}))
```

### Publishing with Synchronization (Flush):

#### Allows you to ensure that all previously sent messages have been acknowledged by the server before proceeding with the program execution.

```go
err := nc.Publish("subject", []byte("Message to be flushed"))
if err != nil {
    log.Fatal(err)
}

err = nc.Flush()
if err != nil {
    log.Fatal(err)
}
log.Println("All messages flushed!")
```

### Monitoring and Statistics:

#### `nats.Conn` provides methods to get statistics on sent and received messages, errors, and latency.

```go
stats := nc.Stats()
log.Printf("Incoming Msgs: %d", stats.InMsgs)
```

### Automatic Reconnection:

#### Support for automatic reconnection to the server when the connection is lost.

```go
nc, err := nats.Connect("nats://localhost:4222", nats.MaxReconnects(10))
```
