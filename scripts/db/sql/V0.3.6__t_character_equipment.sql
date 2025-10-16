create table character_equipments (
    character_id int not null,
    item_instance_id int,
    equip_slot_id int not null,
    
    foreign key (character_id) references characters (character_id),
    -- foreign key (item_instance_id) references item_instances (item_instance_id),
    foreign key (equip_slot_id) references equip_slots (id)
);

create index idx_character_equipments_character_id on character_equipments (character_id);
-- create index idx_character_equipments_item_instance_id on character_equipments (item_instance_id);
create index idx_character_equipments_equip_slot_id on character_equipments (equip_slot_id);