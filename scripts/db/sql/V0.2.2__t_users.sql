CREATE TABLE users (
  id SERIAL,
  username VARCHAR(255) NOT NULL unique,
  password TEXT DEFAULT NULL,
  email VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
);

create unique index idx_users_id on users (id);
CREATE UNIQUE INDEX idx_users_username ON users (username);