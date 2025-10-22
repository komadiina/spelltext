create table item_types (
  id SERIAL,
  code VARCHAR(32) NOT NULL unique,
  name VARCHAR(64) NOT NULL,
  PRIMARY KEY(id)
);

create unique index idx_item_types_code on item_types (code);
-- alter table item_types replica identity full;