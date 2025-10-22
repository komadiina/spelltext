create table item_templates (
  id SERIAL,
  name varchar(255) not null,
  item_type_id int not null,
  equip_slot_id int not null,
  description varchar(255) not null default '',
  gold_price int not null default 0,
  buyable_with_tokens smallint not null default 0,
  token_price smallint not null default 0,
  metadata json default null,

  foreign key (item_type_id) references item_types (id),
  foreign key (equip_slot_id) references equip_slots (id),
  primary key(id)
);

create index idx_item_templates_item_type_id on item_templates (item_type_id);
create index idx_item_templates_equip_slot_id on item_templates (equip_slot_id);

-- alter table item_templates replica identity full;