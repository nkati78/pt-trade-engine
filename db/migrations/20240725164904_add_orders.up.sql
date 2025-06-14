CREATE table orders
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    price      int          NOT NULL,
    quantity   int          NOT NULL,
    side       VARCHAR(255) NOT NULL,
    type       VARCHAR(255) NOT NULL,
    symbol     VARCHAR(255) NOT NULL,
    user_id    UUID         NOT NULL,
    status     VARCHAR(255) NOT NULL,

    created_at TIMESTAMP        DEFAULT current_timestamp,
    updated_at TIMESTAMP        DEFAULT current_timestamp,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX orders_user_id_index ON orders (user_id);
CREATE INDEX orders_symbol_index ON orders (symbol);

CREATE TABLE market_prices (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol          VARCHAR(255) NOT NULL,
    price           INT NOT NULL,
    starting_price  INT,
    today_high     INT,
    today_low      INT,
    today_open    INT,
    yesterday_open INT,
    yesterday_close INT,
    yesterday_high INT,
    yesterday_low  INT,
    trade_date   VARCHAR,
    updated_at      TIMESTAMP    DEFAULT now(),
    reference_id    VARCHAR(255)
);

CREATE INDEX market_prices_symbol_index ON market_prices (symbol);
CREATE UNIQUE INDEX market_prices_symbol_unique_idx ON market_prices (symbol);

INSERT 
    into 
orders 
    (id, price, quantity, side, type, symbol, user_id, status, created_at) 
values 
    ('00dbe482-1642-4c44-a518-f9a4d0df8a44', 20111, 10, 'buy', 'market', 'AAPL', 'f9db6ee0-957d-420e-b3a6-e52613cb63c5', 'filled', '2025-05-11 02:14:42.531863+00:00');

INSERT
into
    orders
(id, price, quantity, side, type, symbol, user_id, status, created_at)
values
    ('00dbe482-1642-4c44-a518-f9a4d0df8a23', 20566, 12, 'buy', 'market', 'IBM', 'f9db6ee0-957d-420e-b3a6-e52613cb63c5', 'filled', '2025-05-10 02:14:42.531863+00:00');

INSERT
into
    orders
(id, price, quantity, side, type, symbol, user_id, status, created_at)
values
    ('00dbe482-1642-4c44-a518-f9a4d0df8a21', 17264, 29, 'sell', 'market', 'GOOGL', 'f9db6ee0-957d-420e-b3a6-e52613cb63c5', 'closed', '2025-05-19 02:14:42.531863+00:00');