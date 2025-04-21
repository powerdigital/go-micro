CREATE TABLE users
(
    id    int primary key auto_increment,
    name  varchar(255),
    email varchar(255),
    phone varchar(255),
    age   tinyint
);

CREATE INDEX idx_users_name ON users (name);
CREATE INDEX idx_users_age ON users (age);
CREATE UNIQUE INDEX idx_users_email_unique ON users (email);
CREATE UNIQUE INDEX idx_users_phone_unique ON users (phone);
