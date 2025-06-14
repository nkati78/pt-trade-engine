CREATE TABLE users
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email      VARCHAR(255) NOT NULL,
    username   VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NOT NULL,
    password_hash TEXT,
    connection_type VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255),

    created_at TIMESTAMP        DEFAULT now(),
    updated_at TIMESTAMP        DEFAULT now()
);

CREATE INDEX users_email_index ON users (email);

INSERT 
    into 
users 
    (id, email, username, first_name, last_name, password_hash, connection_type, provider_id) 
values 
    ('f9db6ee0-957d-420e-b3a6-e52613cb63c5', 'ron.hanks@gumps.com', 'noscopetrader', 'Ronald', 'Hands', '$2a$12$83fXePJg0OMSm8F9uEgK8udQWQ1y.1mVX2pRXf2lwfDdfgwuDL4Yq', 'username-password', ''); 
-- password is ronald