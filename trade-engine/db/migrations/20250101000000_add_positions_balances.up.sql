CREATE table positions
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID         NOT NULL,
    symbol     VARCHAR(255) NOT NULL,
    quantity   int          NOT NULL,
    direction VARCHAR(255) NOT NULL,
    average_price int NOT NULL,
    profit_loss int NOT NULL,
    status    VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    DEFAULT now(),
    updated_at TIMESTAMP    DEFAULT now(),
    order_id UUID NOT NULL
);

CREATE INDEX positions_user_id_index ON positions (user_id);
CREATE INDEX positions_symbol_index ON positions (symbol);
CREATE INDEX positions_order_id_index ON positions (order_id);

CREATE TABLE balances
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID         NOT NULL,
    balance     VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    DEFAULT now(),
    updated_at TIMESTAMP    DEFAULT now()
);

CREATE INDEX balances_user_id_index ON balances (user_id);

INSERT into balances (user_id, balance) values ('f9db6ee0-957d-420e-b3a6-e52613cb63c5', '1000000');

insert into positions 
    (user_id, symbol, quantity, direction, average_price, profit_loss, status, order_id) 
values 
    ('f9db6ee0-957d-420e-b3a6-e52613cb63c5', 'AAPL', 10, 'B', 19972, 15000, 'open', '00dbe482-1642-4c44-a518-f9a4d0df8a44');

CREATE TABLE symbols (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol     VARCHAR(255) NOT NULL,
    exchange   VARCHAR(255) NOT NULL,
    last_trade_price INT,
    last_trade_timestamp TIMESTAMP,
    created_at TIMESTAMP    DEFAULT now(),
    updated_at TIMESTAMP    DEFAULT now()
);

CREATE INDEX symbols_symbol_index ON symbols (symbol);

INSERT into symbols (symbol, exchange) values ('AAPL', 'NASDAQ');
INSERT into symbols (symbol, exchange) values ('GOOGL', 'NASDAQ');
INSERT into symbols (symbol, exchange) values ('MSFT', 'NASDAQ');
INSERT into symbols (symbol, exchange) values ('NFLX', 'NASDAQ');
INSERT into symbols (symbol, exchange) values ('NVDA', 'NASDAQ');
INSERT into symbols (symbol, exchange) values ('IBM', 'NASDAQ');

CREATE TABLE watchlist (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID         NOT NULL,
    symbol     VARCHAR(255) NOT NULL,
    sequence_number  int DEFAULT 0,
    created_at TIMESTAMP    DEFAULT now(),
    updated_at TIMESTAMP    DEFAULT now()
);

CREATE INDEX watchlist_user_id_index ON watchlist (user_id);
CREATE INDEX watchlist_symbol_index ON watchlist (symbol);

INSERT into watchlist (user_id, symbol, sequence_number) values ('f9db6ee0-957d-420e-b3a6-e52613cb63c5', 'AAPL', 1);
INSERT into watchlist (user_id, symbol, sequence_number) values ('f9db6ee0-957d-420e-b3a6-e52613cb63c5', 'GOOGL', 2);
INSERT into watchlist (user_id, symbol, sequence_number) values ('f9db6ee0-957d-420e-b3a6-e52613cb63c5', 'MSFT', 3);
