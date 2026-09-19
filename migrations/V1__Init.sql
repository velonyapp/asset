CREATE TABLE images (
    id          CHAR(36) PRIMARY KEY,
    storage_key VARCHAR(128) NOT NULL UNIQUE,
    ready       BOOLEAN NOT NULL,
    create_time TIMESTAMP(6) NOT NULL,
    update_time TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;
