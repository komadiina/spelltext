create table character_inventory_item_instance (
  character_id int not null,
  item_instance_id int not null,
  foreign key (character_id) references characters (character_id),
  foreign key (item_instance_id) references item_instances (item_instance_id)
);