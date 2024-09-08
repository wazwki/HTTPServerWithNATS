# **How to Set Up NATS in a Golang HTTP Server**

# 1. Installation and Importing the NATS Library:

```bash
go get github.com/nats-io/nats.go
```

```go
import (
    "github.com/nats-io/nats.go"
)
```

# 2. Basic Concepts

## 2.1. Connection:

### Connecting to the Server: An `nats.Conn` object is created, representing the connection to the NATS server.

```go
nc, err := nats.Connect("nats://localhost:4222")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```

## 2.2. Subscription:

### 2.2.1. Standard Subscription: Allows subscribing to specific subjects and receiving messages.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Received message: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}
```

### 2.2.2. Channel Subscription: Messages are placed in a channel, simplifying processing in multi-threaded applications.

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

### 2.2.3. Queue Subscription: Used for load balancing, where multiple subscribers receive messages from a single queue.

```go
sub, err := nc.QueueSubscribe("subject", "worker-group", func(m *nats.Msg) {
    log.Printf("Worker received: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}
```

## 2.3. Publishing:

### 2.3.1. Simple Publishing: Sends a message to a specified subject.

```go
err := nc.Publish("subject", []byte("Hello, NATS!"))
if err != nil {
    log.Fatal(err)
}
```

### 2.3.2. Publishing with Acknowledgment (Ack): Waits for a delivery confirmation from the server.

```go
nc.Publish("subject", []byte("Hello, NATS!"))
err = nc.Flush()
if err != nil {
    log.Fatal(err)
}
log.Println("Message delivered!")
```

## 2.4. Request-Response:

### 2.4.1. Request: Sends a message and waits for a response within a specified time.

```go
msg, err := nc.Request("subject", []byte("Can I get a response?"), 2*time.Second)
if err != nil {
    log.Fatal(err)
}
log.Printf("Received reply: %s", string(msg.Data))
```

### 2.4.2. Response: Handles requests and sends a response back to the client.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Received request: %s", string(m.Data))
    m.Respond([]byte("Sure, here is your response!"))
})
if err != nil {
    log.Fatal(err)
}
```

## 2.5. Auto-Unsubscribe:

### You can set a limit on the number of messages a subscription can receive before it is automatically removed.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Received: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}

sub.AutoUnsubscribe(10)
```

# 3. Advanced Features

## 3.1. Connecting to a Cluster:

### NATS supports connecting to multiple servers, which enhances fault tolerance.

```go
nc, err := nats.Connect("nats://server1:4222,nats://server2:4222")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```

## 3.2. Handling Connection and Error Events:

### `nats.Conn` provides handlers for various events such as connection loss, reconnection, error detection, and connection closure.

```go
nc, err := nats.Connect("nats://localhost:4222",
    nats.DisconnectHandler(func(_ *nats.Conn) {
        log.Println("Disconnected!")
    }),
    nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Printf("Reconnected to %v!\\n", nc.ConnectedUrl())
    }),
    nats.ClosedHandler(func(_ *nats.Conn) {
        log.Println("Connection closed!")
    }),
)
```

## 3.3. Durable Subscriptions:

### For NATS Streaming (an extended version of NATS), durable subscriptions are supported, which preserve state between sessions.

```go
import stan "github.com/nats-io/stan.go"

sc, err := stan.Connect("test-cluster", "client-id")
if err != nil {
    log.Fatal(err)
}
defer sc.Close()

sub, err := sc.Subscribe("subject", func(msg *stan.Msg) {
    log.Printf("Received a message: %s\\n", msg)
}, stan.DurableName("my-durable"))
if err != nil {
    log.Fatal(err)
}
```

## 3.4. Streaming Subscriptions:

### Ability to store messages for subscribers who were temporarily unavailable and deliver missed data upon reconnection.

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

## 3.5. Authorization and Encryption:

### Support for authentication using tokens, passwords, and certificates.

### TLS connections for secure data transmission.

```go
nc, err := nats.Connect("nats://localhost:4222", nats.Secure(&tls.Config{...}))
```

## 3.6. Publishing with Flush:

### Ensures that all previously sent messages have been received by the server before continuing the program execution.

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

## 3.7. Monitoring and Statistics:

### `nats.Conn` provides methods to get statistics on sent and received messages, errors, and latency.

```go
stats := nc.Stats()
log.Printf("Incoming Msgs: %d", stats.InMsgs)
```

## 3.8. Automatic Reconnection:

### Supports automatic reconnection to the server upon disconnection.

```go
nc, err := nats.Connect("nats://localhost:4222", nats.MaxReconnects(10))
```