USE forex;

-- DROP TABLE demo.people IF EXISTS;

CREATE TABLE customer (
    id VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    tier int,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE forex_rate_booking (
    id VARCHAR(50) NOT NULL,
    timestamp timestamp,
    base_currency VARCHAR(3),
    counter_currency VARCHAR(3),
    rate decimal(5,4),
    trade_action VARCHAR(50),
    base_currency_amount decimal(14,4),
    expiry_time timestamp,
    booking_ref VARCHAR(50),
    customer_id VARCHAR(50),
    PRIMARY KEY (id)
);