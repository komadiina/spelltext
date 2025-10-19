create table character_inventories (
  character_id int not null,
  gold int not null,
  boss_tokens int not null,
  base_inventory_size int not null,
  expanded_inventory_size int not null,
  foreign key (character_id) references characters (character_id) on delete cascade,
  primary key(character_id)
);

create unique index idx_character_inventories_id on character_inventories (character_id);

-- alter table character_inventories replica identity full;