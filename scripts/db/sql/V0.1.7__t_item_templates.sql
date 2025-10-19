create table item_templates (
  id SERIAL,
  name varchar(255) not null,
  item_type_id int not null,
  equip_slot_id int not null,
  description text,
  gold_price int not null default 0,
  buyable_with_tokens smallint not null default 0,
  token_price smallint not null default 0,
  metadata json default null,

  foreign key (item_type_id) references item_types (id),
  foreign key (equip_slot_id) references equip_slots (id),
  primary key(id)
);

create unique index idx_item_templates_id on item_templates (id);
create index idx_item_templates_buyable_with_tokens on item_templates (buyable_with_tokens);

-- alter table item_templates replica identity full;