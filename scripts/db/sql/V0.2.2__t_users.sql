CREATE TABLE users (
  id SERIAL,
  username VARCHAR(255) NOT NULL unique,
  password TEXT DEFAULT NULL,
  email VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
);