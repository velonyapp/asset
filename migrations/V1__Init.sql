CREATE TABLE images (
    id          CHAR(36) PRIMARY KEY,
    storage_key VARCHAR(1024) CHARACTER SET ascii NOT NULL UNIQUE,
    status      TEXT NOT NULL,
    create_time TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;