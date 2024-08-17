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

## Основные концепции

### Соединение (Connection):

#### Подключение к серверу: Создается объект nats.Conn, который представляет соединение с NATS-сервером.

```go
nc, err := nats.Connect("nats://localhost:4222")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```


### Подписка (Subscription):

#### Стандартная подписка: Позволяет подписаться на определенные темы и получать сообщения.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Получено сообщение: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}
```

#### Канальная подписка: Сообщения помещаются в канал, что упрощает их обработку в многопоточных приложениях.

```go
ch := make(chan *nats.Msg, 64)
sub, err := nc.ChanSubscribe("subject", ch)
if err != nil {
    log.Fatal(err)
}

go func() {
    for msg := range ch {
        log.Printf("Получено сообщение через канал: %s", string(msg.Data))
    }
}()
```

#### Очередная подписка (Queue Subscription): Используется для реализации балансировки нагрузки, когда несколько подписчиков получают сообщения из одной очереди.

```go
sub, err := nc.QueueSubscribe("subject", "worker-group", func(m *nats.Msg) {
    log.Printf("Worker received: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}
```


### Публикация (Publishing):

#### Простая публикация: Отправляет сообщение в указанную тему.

```go
err := nc.Publish("subject", []byte("Hello, NATS!"))
if err != nil {
    log.Fatal(err)
}
```

#### Публикация с подтверждением (Ack): Ожидание подтверждения о доставке сообщения от сервера.

```go
nc.Publish("subject", []byte("Hello, NATS!"))
err = nc.Flush()
if err != nil {
    log.Fatal(err)
}
log.Println("Сообщение доставлено!")
```


### Запрос-ответ (Request-Reply):

#### Запрос: Отправляет сообщение и ожидает ответа в течение заданного времени.

```go
msg, err := nc.Request("subject", []byte("Can I get a response?"), 2*time.Second)
if err != nil {
    log.Fatal(err)
}
log.Printf("Получен ответ: %s", string(msg.Data))
```

#### Ответ: Обрабатывает запросы и отправляет ответ обратно клиенту.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Получен запрос: %s", string(m.Data))
    m.Respond([]byte("Sure, here is your response!"))
})
if err != nil {
    log.Fatal(err)
}
```


### Авто-удаление подписки (Auto-Unsubscribe):

#### Можно задать лимит на количество сообщений, которые подписка может получить, после чего она автоматически удаляется.

```go
sub, err := nc.Subscribe("subject", func(m *nats.Msg) {
    log.Printf("Received: %s", string(m.Data))
})
if err != nil {
    log.Fatal(err)
}

sub.AutoUnsubscribe(10)
```


## Расширенные возможности

### Соединение с кластером:

#### NATS поддерживает подключение к нескольким серверам, что повышает отказоустойчивость.

```go
nc, err := nats.Connect("nats://server1:4222,nats://server2:4222")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```


### Обработка событий подключения и ошибок:

#### nats.Conn предоставляет обработчики для различных событий, таких как потеря соединения, восстановление, обнаружение ошибок и закрытие соединения.

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


### Постоянная подписка (Durable Subscription):

#### Для NATS Streaming (расширенная версия NATS) поддерживаются постоянные подписки, которые сохраняют состояние между сессиями.

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


### Потоковые подписки (Streaming Subscriptions):

#### Возможность сохранять сообщения для подписчиков, которые были временно недоступны, и передавать им пропущенные данные при восстановлении связи.

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


### Авторизация и шифрование:

#### Поддержка аутентификации с использованием токенов, паролей и сертификатов.
#### TLS-соединения для обеспечения безопасности передачи данных.

```go
nc, err := nats.Connect("nats://localhost:4222", nats.Secure(&tls.Config{...}))
```


### Публикация с синхронизацией (Flush):

#### Позволяет убедиться, что все ранее отправленные сообщения были приняты сервером до продолжения выполнения программы.

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


### Мониторинг и статистика:

#### nats.Conn предоставляет методы для получения статистики по отправленным и полученным сообщениям, ошибкам и времени задержки.

```go
stats := nc.Stats()
log.Printf("Incoming Msgs: %d", stats.InMsgs)
```


### Автоматическое переподключение:

#### Поддержка автоматического переподключения к серверу при разрыве соединения.

```go
nc, err := nats.Connect("nats://localhost:4222", nats.MaxReconnects(10))
```
