CREATE TABLE users
(
    id    serial,
    name  varchar(255),
    email varchar(255),
    phone varchar(255),
    age   smallint
);

CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_age ON users(age);
CREATE UNIQUE INDEX idx_users_email_unique ON users(email);
CREATE UNIQUE INDEX idx_users_phone_unique ON users(phone);
