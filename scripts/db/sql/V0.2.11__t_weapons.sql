create table weapons (
  item_template_id SERIAL,
  prefix varchar(64),
  suffix varchar(64),
  base_damage int not null,
  strength int not null,
  spellpower int not null,
  
  PRIMARY KEY (item_template_id),
  FOREIGN KEY (item_template_id) REFERENCES item_templates (id)
);