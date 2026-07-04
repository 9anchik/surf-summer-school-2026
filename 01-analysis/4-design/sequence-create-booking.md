# Sequence-диаграмма API-взаимодействия

> Этап 3. Проектирование. Как клиент и сервер обмениваются вызовами в критичных сценариях
> бронирования. Контракты API — в многофайловой спецификации
> [api/redocly.yaml](../api/redocly.yaml) (домены `bookings`, `slots`, `auth`).
> Операции: `createBooking`.

> **Сквозные правила взаимодействия.**
> - Все вызовы — с `Authorization: Bearer <token>` (`bearerAuth`); при истёкшем/неверном токене
>   сервер отвечает `401`, клиент уходит на вход [SCR-01](../3-design-brief/scr-01-registration.md).
> - Сервер — **источник истины** по времени и доступности: `slot.start_at` в UTC, тип отмены и
>   наличие мест/экипировки проверяет сервер, клиент их не пересчитывает (R-005, R-021).
> - Запись/отмена **атомарны**: овербукинг и двойная бронь исключены (NFR-8, NFR-9).
> - Таймаут запроса ~10 с; мутации офлайн запрещены — см. единый паттерн Error/Retry (R-020).

## Сценарий 1: Создание брони (`createBooking`, UC-1)

Поток: [SCR-04 «Оформление брони»](../3-design-brief/scr-04-booking.md) → `POST /bookings`
→ [BS-02 «Успешное бронирование»](../3-design-brief/bs-02-booking-success.md). Клиент отправляет
`slot_id`, `seats_count` (1..3) и `rental_count` (0..seats_count). Итоговую цену `price_total`
(RUB, read-only) считает сервер — клиент её не вычисляет, а показывает (R-005, R-010).

```mermaid
sequenceDiagram
    actor User as Клиент
    participant App as Приложение
    participant API as API (bookings)

    Note over App: SCR-04: выбраны места/экипировка,<br/>цена показана из price_total слота
    User->>App: Тап «Записаться»
    App->>App: Генерирует Idempotency-Key (UUID)

    App->>API: POST /bookings<br/>{slot_id, seats_count, rental_count}<br/>Authorization: Bearer, Idempotency-Key
    Note over API: Атомарно: проверка свободных мест/<br/>комплектов экипировки, фиксация цены, списание (NFR-8/9)

    alt Успех
        API-->>App: 201 Booking {id, status: active,<br/>price_total, created_at, slot}
        App-->>User: BS-02 «Запись оформлена» + сводка<br/>(после первой записи — запрос push)
    else Нет свободных мест/экипировки или двойная бронь (409 Conflict)
        API-->>App: 409 {code: slot_full / double_booking,<br/>available_seats, available_rental_equipment}
        App-->>User: Сообщение о нехватке мест/экипировки,<br/>обновление доступности слота
    else Слот отменён центром (410 Gone)
        API-->>App: 410 {code: slot_cancelled, reason}
        App-->>User: «Заезд отменён: <причина>», запись недоступна
    else Невалидные данные (400 / 422)
        API-->>App: 400 BadRequest / 422 Unprocessable
        App-->>User: Подсказка по полям / правилу (макс. 3 места)
    else Токен истёк (401)
        API-->>App: 401 Unauthorized
        App-->>User: Переход на вход (SCR-01)
    else Сеть/сервер/таймаут (~10 c, 5xx)
        API-->>App: Ошибка / нет ответа
        App-->>User: Error + «Повторить» (повтор с тем же<br/>Idempotency-Key — без двойной брони)
    end
```