create table vendors (
  id serial primary key,
  name varchar(64),
  ware_shorthand varchar(64)
);


-- alter table vendors replica identity full;