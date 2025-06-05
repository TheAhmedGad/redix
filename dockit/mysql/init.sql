CREATE DATABASE redix;

USE redix;

CREATE TABLE clients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

INSERT INTO clients (token, is_active) VALUES ('test', TRUE);