CREATE TABLE consumables (
  id SERIAL,
  consumable_template_id INT NOT NULL,
  prefix varchar(64),
  suffix varchar(64),
  uses smallint not null,
  health INT NOT NULL,
  power INT NOT NULL,
  strength INT NOT NULL,
  spellpower INT NOT NULL,
  bonus_damage int not null,
  PRIMARY KEY (id),
  FOREIGN KEY (consumable_template_id) REFERENCES consumable_templates (id)
);


-- alter table consumables replica identity full;