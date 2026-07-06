# Test Cases

## 1. Document Information

| Field | Value |
|-------|-------|
| Project | Apex Karting Booking Service |
| Document | Test Cases |
| Version | 1.0 |
| Status | Draft |
| Author | Daniil Lepetyuha |
| Scope | Backend REST API |
| Target | MVP |

---

# 2. Purpose

This document describes manual test cases for validating the functionality of the Apex Karting Booking Service.

The purpose of testing is to verify that the implemented REST API complies with the functional requirements, business rules, and user scenarios defined during the analysis stage.

The document covers only the MVP functionality implemented in the backend service.

---

# 3. Scope

The following subsystems are covered by testing:

- Authentication
- Slots
- Profile
- Bookings

---

# 4. Test Environment

## Software

| Component | Version |
|-----------|----------|
| Go | 1.25 |
| PostgreSQL | 17 |
| Docker | Latest |
| Docker Compose | Latest |

## Tools

- Postman / Bruno
- Docker Desktop
- psql

---

# 5. Preconditions

Before executing the test cases:

- PostgreSQL database is running.
- Database migrations are successfully applied.
- Backend application is running.
- Test data is inserted into the database.
- User possesses a valid JWT token when testing protected endpoints.
- REST client is configured.

---

# 6. Priority Levels

| Priority | Description |
|----------|-------------|
| High | Critical business functionality |
| Medium | Core application functionality |
| Low | Additional validation |

---

# 7. Authentication

---

## TC-AUTH-001

### Module

Authentication

### Title

Send OTP to a new phone number.

### Priority

High

### Preconditions

- Phone number does not exist in the system.

### Request

```
POST /api/v1/auth/otp/send
```

```json
{
  "phone": "+79991112233"
}
```

### Steps

1. Send POST request.
2. Pass a valid phone number.

### Expected Result

- HTTP Status **200 OK**.
- OTP code is generated.
- OTP expiration time is stored.
- Response contains generated OTP (development mode).

---

## TC-AUTH-002

### Module

Authentication

### Title

Authenticate a new user using a valid OTP.

### Priority

High

### Preconditions

- OTP has been generated.
- OTP has not expired.

### Request

```
POST /api/v1/auth/otp/verify
```

```json
{
  "phone": "+79991112233",
  "code": "123456",
  "name": "Daniil"
}
```

### Steps

1. Send verification request.
2. Provide correct OTP.

### Expected Result

- HTTP Status **200 OK**.
- New user is created.
- JWT Access Token is generated.
- Response contains user identifier and token.

---

## TC-AUTH-003

### Module

Authentication

### Title

Authenticate an existing user.

### Priority

High

### Preconditions

- User already exists.
- OTP has been generated.

### Steps

1. Send verification request.
2. Pass valid OTP.

### Expected Result

- HTTP Status **200 OK**.
- Existing user is returned.
- JWT token is generated.
- Duplicate user is NOT created.

---

## TC-AUTH-004

### Module

Authentication

### Title

Verify invalid OTP.

### Priority

High

### Preconditions

- OTP exists.

### Request

```json
{
  "phone": "+79991112233",
  "code": "000000"
}
```

### Steps

1. Submit incorrect OTP.

### Expected Result

- HTTP Status **401 Unauthorized**.
- User is not authenticated.
- JWT token is not generated.

---

## TC-AUTH-005

### Module

Authentication

### Title

Verify expired OTP.

### Priority

High

### Preconditions

- OTP expiration time has passed.

### Steps

1. Submit expired OTP.

### Expected Result

- HTTP Status **401 Unauthorized**.
- Authentication is rejected.

---

## TC-AUTH-006

### Module

Authentication

### Title

Logout.

### Priority

Medium

### Preconditions

- User is authenticated.

### Request

```
POST /api/v1/auth/logout
```

### Steps

1. Send logout request.
2. Pass Authorization header.

### Expected Result

- HTTP Status **200 OK**.
- Client receives successful logout response.
- Client removes stored JWT token.

---

# 8. Slots

---

## TC-SLOTS-001

### Module

Slots

### Title

Get list of available slots.

### Priority

High

### Preconditions

- Database contains active slots.

### Request

```
GET /api/v1/slots
```

### Steps

1. Send GET request.

### Expected Result

- HTTP Status **200 OK**.
- Response contains slot collection.
- Every slot contains required information.

---

## TC-SLOTS-002

### Module

Slots

### Title

Filter slots by track configuration.

### Priority

Medium

### Preconditions

- Multiple track configurations exist.

### Request

```
GET /api/v1/slots?track_config=short
```

### Steps

1. Send request with filter.

### Expected Result

- HTTP Status **200 OK**.
- Only slots with selected configuration are returned.

---

## TC-SLOTS-003

### Module

Slots

### Title

Return only available slots.

### Priority

Medium

### Preconditions

- Database contains fully booked and available slots.

### Request

```
GET /api/v1/slots?only_available=true
```

### Steps

1. Send filtered request.

### Expected Result

- HTTP Status **200 OK**.
- Fully booked slots are excluded.

---

## TC-SLOTS-004

### Module

Slots

### Title

Get slot details.

### Priority

High

### Preconditions

- Slot exists.

### Request

```
GET /api/v1/slots/{id}
```

### Steps

1. Request slot by identifier.

### Expected Result

- HTTP Status **200 OK**.
- Complete slot information is returned.

---

## TC-SLOTS-005

### Module

Slots

### Title

Request non-existing slot.

### Priority

Medium

### Preconditions

- Slot identifier does not exist.

### Steps

1. Request slot using invalid identifier.

### Expected Result

- HTTP Status **404 Not Found**.
- Error response is returned.

---

## TC-SLOTS-006

### Module

Slots

### Title

Verify free seats calculation.

### Priority

Medium

### Preconditions

- Slot contains partially booked seats.

### Steps

1. Request slot details.

### Expected Result

- Number of free seats equals:

```
total_seats - booked_seats
```

---

## TC-SLOTS-007

### Module

Slots

### Title

Verify available rental equipment count.

### Priority

Medium

### Preconditions

- Slot contains booked rental equipment.

### Steps

1. Request slot details.

### Expected Result

- Available rental equipment equals:

```
rental_equipment_total -
rental_equipment_booked
```

---

## TC-SLOTS-008

### Module

Slots

### Title

Verify slot sorting.

### Priority

Low

### Preconditions

- Database contains several future slots.

### Steps

1. Request slot list.

### Expected Result

- Slots are sorted by start date in ascending order.

---

## TC-SLOTS-009

### Module

Slots

### Title

Request slots from empty database.

### Priority

Low

### Preconditions

- Slots table contains no active records.

### Steps

1. Request slot list.

### Expected Result

- HTTP Status **200 OK**.
- Empty array is returned.

---

## TC-SLOTS-010

### Module

Slots

### Title

Verify response model.

### Priority

Low

### Preconditions

- Slot exists.

### Steps

1. Request slot.

### Expected Result

Response contains:

- id
- start_at
- arrival_at
- track configuration
- marshal
- address
- meeting point
- prices
- free seats
- rental equipment
- slot status

# 9. Profile

---

## TC-PROFILE-001

### Module

Profile

### Title

Get current user profile.

### Priority

High

### Preconditions

- User is authenticated.
- JWT Access Token is valid.

### Request

```
GET /api/v1/profile
```

### Steps

1. Send request with Authorization header.

### Expected Result

- HTTP Status **200 OK**.
- User profile is returned.
- Response contains:
  - id
  - name
  - phone
  - created_at
  - updated_at

---

## TC-PROFILE-002

### Module

Profile

### Title

Update user name.

### Priority

High

### Preconditions

- User is authenticated.

### Request

```
PATCH /api/v1/profile
```

```json
{
  "name": "Daniil Lepetyuha"
}
```

### Steps

1. Send PATCH request.
2. Pass Authorization header.

### Expected Result

- HTTP Status **200 OK**.
- Name is successfully updated.
- updated_at value is changed.

---

## TC-PROFILE-003

### Module

Profile

### Title

Update phone number.

### Priority

Medium

### Preconditions

- User is authenticated.

### Request

```json
{
  "phone": "+79995554433"
}
```

### Steps

1. Send PATCH request.

### Expected Result

- HTTP Status **200 OK**.
- Phone number is updated.

---

## TC-PROFILE-004

### Module

Profile

### Title

Update profile with empty name.

### Priority

Medium

### Preconditions

- User is authenticated.

### Request

```json
{
    "name":""
}
```

### Steps

1. Send PATCH request.

### Expected Result

- HTTP Status **400 Bad Request**.
- Validation error is returned.
- Profile remains unchanged.

---

## TC-PROFILE-005

### Module

Profile

### Title

Delete profile.

### Priority

High

### Preconditions

- User is authenticated.
- User has active bookings.

### Request

```
DELETE /api/v1/profile
```

### Steps

1. Send DELETE request.
2. Pass Authorization header.

### Expected Result

- HTTP Status **200 OK**.
- User is marked as deleted.
- Active bookings are cancelled.
- Slot counters are restored.

---

# 10. Bookings

---

## TC-BOOK-001

### Module

Bookings

### Title

Create booking for one participant.

### Priority

High

### Preconditions

- User is authenticated.
- Slot exists.
- Slot has free seats.

### Request

```
POST /api/v1/bookings
```

```json
{
  "slot_id":"slot-id",
  "equipment":[
    "own"
  ]
}
```

### Steps

1. Send request.
2. Pass Authorization header.
3. Pass Idempotency-Key header.

### Expected Result

- HTTP Status **201 Created**.
- Booking is created.
- Slot counters are updated.

---

## TC-BOOK-002

### Module

Bookings

### Title

Create booking for three participants.

### Priority

High

### Preconditions

- Slot has at least three free seats.

### Request

```json
{
  "slot_id":"slot-id",
  "equipment":[
    "own",
    "own",
    "rental"
  ]
}
```

### Steps

1. Send booking request.

### Expected Result

- HTTP Status **201 Created**.
- Booking contains three participants.
- Rental counter is increased by one.

---

## TC-BOOK-003

### Module

Bookings

### Title

Create booking when slot is full.

### Priority

High

### Preconditions

- Slot has no available seats.

### Steps

1. Send booking request.

### Expected Result

- HTTP Status **409 Conflict**.
- Booking is not created.

---

## TC-BOOK-004

### Module

Bookings

### Title

Create booking without available rental equipment.

### Priority

High

### Preconditions

- Rental equipment is fully booked.

### Steps

1. Send booking request requesting rental equipment.

### Expected Result

- HTTP Status **409 Conflict**.
- Booking is rejected.

---

## TC-BOOK-005

### Module

Bookings

### Title

Repeat request using the same Idempotency-Key.

### Priority

High

### Preconditions

- Booking has already been created.

### Steps

1. Repeat identical request.
2. Use the same Idempotency-Key.

### Expected Result

- Existing booking is returned.
- Duplicate booking is not created.

---

## TC-BOOK-006

### Module

Bookings

### Title

View booking list.

### Priority

High

### Preconditions

- User has bookings.

### Request

```
GET /api/v1/bookings
```

### Steps

1. Send request.

### Expected Result

- HTTP Status **200 OK**.
- User booking list is returned.

---

## TC-BOOK-007

### Module

Bookings

### Title

View booking details.

### Priority

High

### Preconditions

- Booking exists.

### Request

```
GET /api/v1/bookings/{id}
```

### Steps

1. Request booking by identifier.

### Expected Result

- HTTP Status **200 OK**.
- Booking information is returned.
- Equipment list is returned.

---

## TC-BOOK-008

### Module

Bookings

### Title

Cancel booking more than three hours before slot start.

### Priority

High

### Preconditions

- Booking is active.
- Slot starts in more than three hours.

### Request

```
POST /api/v1/bookings/{id}/cancel
```

### Steps

1. Send cancellation request.

### Expected Result

- HTTP Status **200 OK**.
- Booking status becomes **cancelled**.
- Slot counters are restored.

---

## TC-BOOK-009

### Module

Bookings

### Title

Late cancellation.

### Priority

High

### Preconditions

- Booking is active.
- Less than three hours remain before slot start.

### Steps

1. Send cancellation request.

### Expected Result

- HTTP Status **200 OK**.
- Booking status becomes **late_cancel**.
- Slot counters remain unchanged.

---

## TC-BOOK-010

### Module

Bookings

### Title

Cancel already cancelled booking.

### Priority

Medium

### Preconditions

- Booking status is not active.

### Steps

1. Send cancellation request again.

### Expected Result

- HTTP Status **409 Conflict**.
- Booking status does not change.

---

## TC-BOOK-011

### Module

Bookings

### Title

Cancel another user's booking.

### Priority

High

### Preconditions

- Booking belongs to another user.

### Steps

1. Send cancellation request.

### Expected Result

- HTTP Status **404 Not Found**.
- Booking is not modified.

---

## TC-BOOK-012

### Module

Bookings

### Title

Create booking after slot has started.

### Priority

High

### Preconditions

- Slot start time has already passed.

### Steps

1. Send booking request.

### Expected Result

- Booking creation is rejected.
- Appropriate error is returned.

---

## TC-BOOK-013

### Module

Bookings

### Title

Verify booking price calculation.

### Priority

Medium

### Preconditions

- Slot exists.
- Rental equipment selected.

### Steps

1. Create booking.

### Expected Result

Total price equals:

```
slot_price × seats_count +
rental_price × rental_count
```

---

## TC-BOOK-014

### Module

Bookings

### Title

Verify booking seat records.

### Priority

Medium

### Preconditions

- Booking created.

### Steps

1. Create booking.
2. Verify database.

### Expected Result

One record exists in **booking_seats** for every participant.

---

## TC-BOOK-015

### Module

Bookings

### Title

Verify booking transaction consistency.

### Priority

High

### Preconditions

- Database is available.

### Steps

1. Simulate database error during booking creation.

### Expected Result

- Transaction is rolled back.
- Booking is not created.
- Slot counters remain unchanged.

# 11. Negative Scenarios

---

## TC-NEG-001

### Module

Authentication

### Title

Access protected endpoint without JWT token.

### Priority

High

### Preconditions

None.

### Request

```
GET /api/v1/profile
```

### Steps

1. Send request without Authorization header.

### Expected Result

- HTTP Status **401 Unauthorized**.
- Access is denied.

---

## TC-NEG-002

### Module

Authentication

### Title

Access protected endpoint using invalid JWT token.

### Priority

High

### Preconditions

Invalid JWT token.

### Steps

1. Send request with malformed JWT.

### Expected Result

- HTTP Status **401 Unauthorized**.
- Access is denied.

---

## TC-NEG-003

### Module

Authentication

### Title

Access protected endpoint using expired JWT token.

### Priority

High

### Preconditions

Expired JWT token.

### Steps

1. Send request with expired JWT.

### Expected Result

- HTTP Status **401 Unauthorized**.
- User is not authenticated.

---

## TC-NEG-004

### Module

Validation

### Title

Request with malformed JSON body.

### Priority

Medium

### Preconditions

Protected endpoint is available.

### Steps

1. Send request with invalid JSON.

Example:

```json
{
  "slot_id":
```

### Expected Result

- HTTP Status **400 Bad Request**.
- Validation error is returned.

---

## TC-NEG-005

### Module

Bookings

### Title

Create booking without Idempotency-Key.

### Priority

High

### Preconditions

User is authenticated.

### Steps

1. Send POST /api/v1/bookings.
2. Do not pass Idempotency-Key header.

### Expected Result

- HTTP Status **400 Bad Request**.
- Booking is not created.

---

## TC-NEG-006

### Module

Bookings

### Title

Request booking with invalid identifier.

### Priority

Medium

### Preconditions

User is authenticated.

### Steps

1. Send GET /api/v1/bookings/{id}.
2. Pass random UUID.

### Expected Result

- HTTP Status **404 Not Found**.

---

## TC-NEG-007

### Module

Profile

### Title

Update profile using invalid request body.

### Priority

Medium

### Preconditions

User is authenticated.

### Steps

1. Send PATCH /profile.
2. Pass unsupported JSON fields.

### Expected Result

- HTTP Status **400 Bad Request**.

---

## TC-NEG-008

### Module

Slots

### Title

Request slot using invalid identifier.

### Priority

Medium

### Preconditions

None.

### Steps

1. Send GET /slots/{id}.
2. Pass non-existing UUID.

### Expected Result

- HTTP Status **404 Not Found**.

---

# 12. Security Test Cases

---

## TC-SEC-001

### Title

User cannot access another user's booking.

### Priority

High

### Preconditions

Two registered users exist.

### Steps

1. Authenticate as User A.
2. Request booking belonging to User B.

### Expected Result

- HTTP Status **404 Not Found**.
- Booking information is not disclosed.

---

## TC-SEC-002

### Title

User cannot cancel another user's booking.

### Priority

High

### Preconditions

Booking belongs to another user.

### Steps

1. Authenticate as another user.
2. Send cancellation request.

### Expected Result

- HTTP Status **404 Not Found**.
- Booking remains unchanged.

---

## TC-SEC-003

### Title

Deleted user cannot access protected resources.

### Priority

Medium

### Preconditions

User profile has been deleted.

### Steps

1. Send GET /profile using previously issued JWT.

### Expected Result

- Request is rejected.
- Deleted profile is not returned.

---

# 13. Traceability Matrix

| Requirement | Covered Test Cases |
|-------------|-------------------|
| User authentication | TC-AUTH-001 ... TC-AUTH-006 |
| View slots | TC-SLOTS-001 ... TC-SLOTS-010 |
| Manage profile | TC-PROFILE-001 ... TC-PROFILE-005 |
| Create booking | TC-BOOK-001 ... TC-BOOK-005 |
| View bookings | TC-BOOK-006 ... TC-BOOK-007 |
| Cancel booking | TC-BOOK-008 ... TC-BOOK-011 |
| Booking calculations | TC-BOOK-013 ... TC-BOOK-015 |
| Validation | TC-NEG-001 ... TC-NEG-008 |
| Security | TC-SEC-001 ... TC-SEC-003 |

---

# 14. Test Coverage Summary

| Module | Number of Test Cases |
|---------|---------------------:|
| Authentication | 6 |
| Slots | 10 |
| Profile | 5 |
| Bookings | 15 |
| Negative Scenarios | 8 |
| Security | 3 |

**Total:** **47 manual test cases**

---

# 15. Test Execution Summary

| Metric | Value |
|--------|------:|
| Total Test Cases | 47 |
| Passed | |
| Failed | |
| Blocked | |
| Not Executed | |

---

# 16. Defect Log

| Defect ID | Test Case | Description | Severity | Status |
|-----------|-----------|-------------|----------|--------|
| | | | | |

---

# 17. Conclusion

The test cases described in this document provide comprehensive manual coverage of the MVP functionality implemented in the Apex Karting Booking Service.

The scenarios validate:

- authentication and authorization;
- slot retrieval and filtering;
- profile management;
- booking creation and cancellation;
- transaction consistency;
- idempotent request processing;
- validation of incorrect requests;
- security restrictions for protected resources.

Together with the implemented unit tests for the service layer, these manual test cases provide confidence that the backend satisfies the business and functional requirements defined for the MVP.