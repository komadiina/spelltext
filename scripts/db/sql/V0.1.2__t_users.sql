CREATE TABLE users (
  id SERIAL,
  username VARCHAR(255) NOT NULL unique,
  password TEXT DEFAULT NULL,
  email VARCHAR(255) NOT NULL,
  selected_character_id int default 0,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_users_username ON users (username);

-- alter table users replica identity full;