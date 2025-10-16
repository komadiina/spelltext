create table item_instances (
  item_instance_id BIGSERIAL,
  item_id INT NOT NULL,
  owner_character_id INT DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  metadata JSON DEFAULT NULL,
  PRIMARY KEY(item_instance_id),
  FOREIGN KEY (item_id) REFERENCES items(id),
  FOREIGN KEY (owner_character_id) REFERENCES characters(character_id)
);

create index idx_item_instances_item_id on item_instances (item_id);
create index idx_item_instances_owner_character_id on item_instances (owner_character_id);