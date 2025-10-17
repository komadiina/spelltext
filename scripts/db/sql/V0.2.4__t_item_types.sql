create table item_types (
  id SERIAL,
  code VARCHAR(32) NOT NULL unique,
  name VARCHAR(64) NOT NULL,
  PRIMARY KEY(id)
);


CREATE UNIQUE INDEX idx_item_types_code ON item_types (code);