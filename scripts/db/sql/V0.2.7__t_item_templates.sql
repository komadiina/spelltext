create table item_templates (
  id SERIAL,
  name varchar(255) not null,
  item_type_id smallint not null,
  rarity smallint  not null default 0,
  stackable smallint not null default 0,
  stack_size integer  not null default 1,
  bind_on_pickup smallint not null default 1,
  description text,
  metadata json default null,
  foreign key (item_type_id) references item_types (id),
  primary key(id)
);