create table item_instances (
  item_instance_id BIGSERIAL,
  template_id INT NOT NULL,
  owner_user_id INT DEFAULT NULL,
  bound_character_id INT DEFAULT NULL,
  durability INT DEFAULT NULL,
  durability_max INT DEFAULT NULL,
  stack_count INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  metadata JSON DEFAULT NULL,
  PRIMARY KEY(item_instance_id),
  FOREIGN KEY (template_id) REFERENCES item_templates(id),
  FOREIGN KEY (owner_user_id) REFERENCES users(id),
  FOREIGN KEY (bound_character_id) REFERENCES characters(character_id)
);