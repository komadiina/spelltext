create table vendors (
  id serial primary key,
  name varchar(64),
  ware_shorthand varchar(64)
);

create unique index idx_vendors_id on vendors (id);