create table character_inventory_item_instance (
  character_id int not null,
  item_instance_id int not null,
  foreign key (character_id) references characters (character_id) on delete cascade,
  foreign key (item_instance_id) references item_instances (item_instance_id) on delete cascade,
  primary key (character_id, item_instance_id)
);

-- alter table character_inventory_item_instance replica identity full;