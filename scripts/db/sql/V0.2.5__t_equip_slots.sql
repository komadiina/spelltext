create table equip_slots (
  id SERIAL,
  code VARCHAR(32) NOT NULL unique,
  name VARCHAR(64) NOT NULL,
  PRIMARY KEY(id)
);

create unique index idx_equip_slots_id on equip_slots (id);
create unique index idx_equip_slots_code on equip_slots (code);