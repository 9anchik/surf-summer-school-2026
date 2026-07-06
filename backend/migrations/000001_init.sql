CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT,
    phone TEXT NOT NULL UNIQUE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE otp_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone TEXT NOT NULL,
    code TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    attempts_left INT NOT NULL DEFAULT 5,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE track_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    description TEXT NOT NULL,
    capacity_cap INT NOT NULL CHECK (capacity_cap > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE marshals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE slots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    track_config_id UUID NOT NULL REFERENCES track_configs(id),
    marshal_id UUID REFERENCES marshals(id),

    start_at TIMESTAMPTZ NOT NULL,
    arrival_at TIMESTAMPTZ NOT NULL,

    total_seats INT NOT NULL CHECK (total_seats > 0),
    booked_seats INT NOT NULL DEFAULT 0 CHECK (booked_seats >= 0),

    rental_equipment_total INT NOT NULL DEFAULT 0 CHECK (rental_equipment_total >= 0),
    rental_equipment_booked INT NOT NULL DEFAULT 0 CHECK (rental_equipment_booked >= 0),

    price INT NOT NULL CHECK (price >= 0),
    rental_price INT NOT NULL DEFAULT 0 CHECK (rental_price >= 0),
    currency TEXT NOT NULL DEFAULT 'RUB',

    address TEXT NOT NULL,
    meeting_point_name TEXT NOT NULL,
    meeting_point_lat NUMERIC(9, 6) NOT NULL,
    meeting_point_lng NUMERIC(9, 6) NOT NULL,

    status TEXT NOT NULL DEFAULT 'active',
    cancel_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (booked_seats <= total_seats),
    CHECK (rental_equipment_booked <= rental_equipment_total),
    CHECK (status IN ('active', 'cancelled_by_center'))
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    slot_id UUID NOT NULL REFERENCES slots(id),

    seats_count INT NOT NULL CHECK (seats_count BETWEEN 1 AND 3),
    rental_count INT NOT NULL DEFAULT 0 CHECK (rental_count >= 0),
    price_total INT NOT NULL CHECK (price_total >= 0),
    currency TEXT NOT NULL DEFAULT 'RUB',

    status TEXT NOT NULL DEFAULT 'active',
    idempotency_key TEXT NOT NULL,

    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (rental_count <= seats_count),
    CHECK (status IN ('active', 'cancelled', 'late_cancel', 'cancelled_by_center')),

    UNIQUE (user_id, idempotency_key)
);

CREATE TABLE booking_seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    equipment_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (equipment_type IN ('own', 'rental'))
);

CREATE INDEX idx_slots_start_at ON slots(start_at);
CREATE INDEX idx_slots_status ON slots(status);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_slot_id ON bookings(slot_id);
CREATE INDEX idx_otp_codes_phone ON otp_codes(phone);

INSERT INTO track_configs (name, code, type, description, capacity_cap)
VALUES
    ('Короткая трасса', 'short', 'beginner', 'Новичковая конфигурация трассы', 8),
    ('Длинная трасса', 'long', 'experienced', 'Опытная конфигурация трассы', 14);

INSERT INTO marshals (name)
VALUES
    ('Алексей'),
    ('Игорь');

INSERT INTO slots (
    track_config_id,
    marshal_id,
    start_at,
    arrival_at,
    total_seats,
    rental_equipment_total,
    price,
    rental_price,
    address,
    meeting_point_name,
    meeting_point_lat,
    meeting_point_lng
)
SELECT
    tc.id,
    m.id,
    NOW() + INTERVAL '1 day',
    NOW() + INTERVAL '1 day' - INTERVAL '15 minutes',
    tc.capacity_cap,
    6,
    2500,
    500,
    'Картинг-центр Apex, Москва',
    'Вход у ресепшена',
    55.751244,
    37.618423
FROM track_configs tc
CROSS JOIN marshals m
LIMIT 4;