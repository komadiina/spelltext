CREATE TABLE consumables (
  item_template_id INT NOT NULL,
  uses smallint  not null,
  prefix varchar(64) not null,
  suffix varchar(64) not null,
  healing INT NOT NULL,
  power INT NOT NULL,
  strength INT NOT NULL,
  spellpower INT NOT NULL,
  PRIMARY KEY (item_template_id),
  FOREIGN KEY (item_template_id) REFERENCES item_templates (id)
);