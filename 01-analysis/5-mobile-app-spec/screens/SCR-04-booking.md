# Оформление брони (Modal / Bottom Sheet)

**ID:** SCR-04  
**Тип:** Bottom Sheet / Модалка  
**Домен:** Mobile App / Booking  
**Приоритет:** Critical  
**Статус:** Актуален  
**Функциональные блоки:** FB-BOOK-001

---

## Обзор

Многошаговая шторка оформления брони: выбор числа мест (1–3), выбор экипировки для каждого места (своя/прокат), показ итоговой цены (`price_total`) из API и подтверждение брони. Idempotency-key обязателен.

### User Story

> Как клиент, хочу указать количество мест и экипировку, чтобы успешно забронировать себя и гостей.

---

## Навигация

Вход: SCR-03 (карточка слота) → Tap «Записаться».  
Исход: Экран успеха в модалке → возврат к SCR-02.

---

## Входные данные

slot_id, free_seats, free_rental_equipment, track_config.capacity_cap, price (серверный price_total в ответе POST)

---

## Применяемые логики

LOGIC-03 (создание брони) — idempotency, валидация лимитов; LOGIC-02 (перепроверка доступности).

---

## Используемые запросы

### POST /bookings

**Body:** resource_id/slot_id, seats_count, rental_count, idempotency_key  
**Обработка:** 201 — success (booking_id), 409 — conflict (show updated slots), 422 — validation error.

---

## Макет и элементы

Stepper количества мест, список участников с тогглами экипировки, summary price, CTA "Подтвердить и забронировать".

---

## Состояния

Loading — spinner in CTA; Success — show BS success (BS-02); Conflict — show updated slots and lead user back to list.

---

## Действия

Подтвердить → POST /bookings с idempotency_key, при успехе показать страницу успеха с summary и опцией добавить в календарь.

---

## Связанные требования

### Функциональные (REQ-FUNC-*)

| ID | Название | Приоритет |
|---|---|---|
| FR-10 | Запись на слот | Must |
| FR-11 | Выбор экипировки | Must |
| FR-12 | Бронирование до 3 мест | Must |
| FR-13 | Ограничение max_seats | Must |
| FR-14 | Учёт прокатного фонда | Must |
| FR-15 | Защита от overbooking / double booking | Must |
| FR-45 | Использовать server price_total | Must |

### Интеграции (REQ-INT-*)

| ID | Название | API operationId (контракт) |
|---|---|---|
| REQ-INT-BOOKINGS | Backend Bookings API | createBooking — POST /bookings — см. /01-analysis/api/bookings/bookings.yaml (использовать заголовок Idempotency-Key)
| REQ-INT-SLOTS | Backend Slots API | (re-check availability) listSlots / getSlot — см. /01-analysis/api/slots/slots.yaml


## Критерии приёмки

- Idempotency предотвращает дубли
- Ограничения на места и прокат соблюдаются
- При нехватке — информативная подсказка и блокировки
