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