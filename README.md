# exponent-server-sdk-golang

Go client for the [Expo Push Notifications API](https://docs.expo.dev/push-notifications/sending-notifications/).

This is a maintained fork of [oliveroneill/exponent-server-sdk-golang](https://github.com/oliveroneill/exponent-server-sdk-golang) with fixes and extensions used in production at Rentifly.

## Features

- Send push notifications via `POST /push/send`
- Fetch delivery receipts via `POST /push/getReceipts`
- Correct handling of multiple recipients in a single `PushMessage`
- Parsing of nested `details` payloads from Expo (`fcm`, `apns`, nested `error` objects)
- Typed errors for common Expo error codes (`DeviceNotRegistered`, `MessageTooBig`, `MessageRateExceeded`)

## Installation

```bash
go get github.com/rentifly/exponent-server-sdk-golang/sdk
```

## Quick start

```go
package main

import (
	"fmt"
	"log"

	expo "github.com/rentifly/exponent-server-sdk-golang/sdk"
)

func main() {
	pushToken, err := expo.NewExponentPushToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]")
	if err != nil {
		log.Fatal(err)
	}

	client := expo.NewPushClient(nil)

	responses, err := client.Publish(&expo.PushMessage{
		To:       []expo.ExponentPushToken{pushToken},
		Title:    "Notification title",
		Body:     "Notification body",
		Data:     map[string]string{"key": "value"},
		Sound:    "default",
		Priority: expo.DefaultPriority,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, response := range responses {
		if err := response.ValidateResponse(); err != nil {
			fmt.Println("failed:", err)
			continue
		}

		fmt.Println("ticket id:", response.ID)
	}
}
```

## Multiple recipients

`Publish` returns one `PushResponse` per recipient in `message.To`, in the same order:

```go
responses, err := client.Publish(&expo.PushMessage{
	To: []expo.ExponentPushToken{token1, token2},
	Body: "Hello",
})
// len(responses) == 2
```

## Delivery receipts

After sending, Expo returns ticket IDs. Wait about 15 minutes (Expo recommends checking receipts later; 1 minute is often enough in practice), then fetch delivery status:

```go
receipts, err := client.GetReceipts([]string{
	"019f0039-9cb1-7088-a76b-504af6cea9b7",
})
if err != nil {
	log.Fatal(err)
}

for ticketID, receipt := range receipts {
	if err := receipt.ValidateReceipt(); err != nil {
		fmt.Println(ticketID, "delivery failed:", err)
		continue
	}

	fmt.Println(ticketID, "delivered")
}
```

Receipt `details` may contain plain strings or nested provider objects. The SDK normalizes them into `PushDetails` (`map[string]string`):

```json
{
  "error": "DeviceNotRegistered",
  "fcm": { "error": "NotRegistered" }
}
```

Access parsed values via `receipt.Details["error"]`, `receipt.Details["fcm"]`, etc.

## Access token

If your Expo project requires an access token:

```go
client := expo.NewPushClient(&expo.ClientConfig{
	AccessToken: "your-expo-access-token",
})
```

## Error handling

### Push tickets (`Publish`)

- `PushServerError` — entire request failed
- `PushResponseError` — single recipient failed
- `DeviceNotRegisteredError` — stop sending to this token
- `MessageTooBigError` — payload exceeds 4096 bytes
- `MessageRateExceededError` — rate limit exceeded, use backoff

### Delivery receipts (`GetReceipts`)

- `ReceiptsServerError` — entire request failed
- `PushReceiptError` — receipt indicates delivery failure
- `DeviceNotRegisteredReceiptError`, `MessageTooBigReceiptError`, `MessageRateExceededReceiptError`

## API reference

| Method | Expo endpoint | Description |
|--------|---------------|-------------|
| `Publish` | `POST /--/api/v2/push/send` | Send one message to one or more tokens |
| `PublishMultiple` | `POST /--/api/v2/push/send` | Send multiple messages in one request |
| `GetReceipts` | `POST /--/api/v2/push/getReceipts` | Fetch delivery receipts by ticket IDs |

Docs: https://docs.expo.dev/push-notifications/sending-notifications/

## Changelog (Rentifly fork)

| Version | Changes |
|---------|---------|
| `v1.2.1` | Parse nested `details` in tickets and receipts (`fcm`, `apns`, nested `error`) |
| `v1.2.0` | Add `GetReceipts` |
| `v1.1.0` | Fix ticket mapping when one message has multiple recipients |

## License

MIT

---

# exponent-server-sdk-golang (русский)

Go-клиент для [Expo Push Notifications API](https://docs.expo.dev/push-notifications/sending-notifications/).

Это поддерживаемый форк [oliveroneill/exponent-server-sdk-golang](https://github.com/oliveroneill/exponent-server-sdk-golang) с исправлениями и доработками, которые используются в production в Rentifly.

## Возможности

- Отправка push-уведомлений через `POST /push/send`
- Получение delivery receipts через `POST /push/getReceipts`
- Корректная обработка нескольких получателей в одном `PushMessage`
- Разбор вложенных `details` от Expo (`fcm`, `apns`, вложенные объекты `error`)
- Типизированные ошибки для распространённых кодов Expo (`DeviceNotRegistered`, `MessageTooBig`, `MessageRateExceeded`)

## Установка

```bash
go get github.com/rentifly/exponent-server-sdk-golang/sdk
```

## Быстрый старт

```go
package main

import (
	"fmt"
	"log"

	expo "github.com/rentifly/exponent-server-sdk-golang/sdk"
)

func main() {
	pushToken, err := expo.NewExponentPushToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]")
	if err != nil {
		log.Fatal(err)
	}

	client := expo.NewPushClient(nil)

	responses, err := client.Publish(&expo.PushMessage{
		To:       []expo.ExponentPushToken{pushToken},
		Title:    "Заголовок",
		Body:     "Текст уведомления",
		Data:     map[string]string{"key": "value"},
		Sound:    "default",
		Priority: expo.DefaultPriority,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, response := range responses {
		if err := response.ValidateResponse(); err != nil {
			fmt.Println("ошибка:", err)
			continue
		}

		fmt.Println("ticket id:", response.ID)
	}
}
```

## Несколько получателей

`Publish` возвращает один `PushResponse` на каждый токен в `message.To`, в том же порядке:

```go
responses, err := client.Publish(&expo.PushMessage{
	To: []expo.ExponentPushToken{token1, token2},
	Body: "Привет",
})
// len(responses) == 2
```

## Delivery receipts

После отправки Expo возвращает ticket ID. Рекомендуется подождать (Expo советует проверять receipts позже; на практике часто хватает ~1 минуты), затем запросить статус доставки:

```go
receipts, err := client.GetReceipts([]string{
	"019f0039-9cb1-7088-a76b-504af6cea9b7",
})
if err != nil {
	log.Fatal(err)
}

for ticketID, receipt := range receipts {
	if err := receipt.ValidateReceipt(); err != nil {
		fmt.Println(ticketID, "доставка не удалась:", err)
		continue
	}

	fmt.Println(ticketID, "доставлено")
}
```

Поле `details` в receipt может содержать как строки, так и вложенные объекты провайдеров. SDK нормализует их в `PushDetails` (`map[string]string`):

```json
{
  "error": "DeviceNotRegistered",
  "fcm": { "error": "NotRegistered" }
}
```

Значения доступны через `receipt.Details["error"]`, `receipt.Details["fcm"]` и т.д.

## Access token

Если для проекта Expo нужен access token:

```go
client := expo.NewPushClient(&expo.ClientConfig{
	AccessToken: "your-expo-access-token",
})
```

## Обработка ошибок

### Push tickets (`Publish`)

- `PushServerError` — упал весь запрос
- `PushResponseError` — ошибка для одного получателя
- `DeviceNotRegisteredError` — перестать слать на этот токен
- `MessageTooBigError` — payload больше 4096 байт
- `MessageRateExceededError` — превышен rate limit, нужен backoff

### Delivery receipts (`GetReceipts`)

- `ReceiptsServerError` — упал весь запрос
- `PushReceiptError` — receipt сообщает об ошибке доставки
- `DeviceNotRegisteredReceiptError`, `MessageTooBigReceiptError`, `MessageRateExceededReceiptError`

## Справка по API

| Метод | Endpoint Expo | Описание |
|-------|---------------|----------|
| `Publish` | `POST /--/api/v2/push/send` | Отправить одно сообщение одному или нескольким токенам |
| `PublishMultiple` | `POST /--/api/v2/push/send` | Отправить несколько сообщений одним запросом |
| `GetReceipts` | `POST /--/api/v2/push/getReceipts` | Получить delivery receipts по ticket ID |

Документация: https://docs.expo.dev/push-notifications/sending-notifications/

## История изменений (форк Rentifly)

| Версия | Изменения |
|--------|-----------|
| `v1.2.1` | Разбор вложенных `details` в tickets и receipts (`fcm`, `apns`, вложенный `error`) |
| `v1.2.0` | Добавлен `GetReceipts` |
| `v1.1.0` | Исправлено сопоставление tickets при нескольких получателях в одном сообщении |

## Лицензия

MIT
