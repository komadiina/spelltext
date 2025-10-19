create table character_equipments (
    character_id int not null,
    equip_slot_id int not null,
    item_instance_id int,
    
    foreign key (character_id) references characters (character_id) on delete cascade,
    foreign key (equip_slot_id) references equip_slots (id),
    primary key (character_id, equip_slot_id)
);

create index idx_character_equipments_character_id on character_equipments (character_id);
create index idx_character_equipments_equip_slot_id on character_equipments (equip_slot_id);

-- alter table character_equipments replica identity full;