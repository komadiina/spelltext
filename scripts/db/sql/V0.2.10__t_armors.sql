CREATE TABLE armors (
  item_template_id INT NOT NULL,
  equip_slot_id smallint not null,
  armor int not null,
  health int not null,
  power int not null,
  strength int not null,
  spellpower int not null,
  foreign key (item_template_id) references item_templates (id) on delete cascade,
  foreign key (equip_slot_id) references equip_slots (id)
);