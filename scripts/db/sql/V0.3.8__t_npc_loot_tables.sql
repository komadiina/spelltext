create table npc_loot_tables (
    npc_template_id int not null,
    item_id int not null,
    chance float not null,
    quantity int not null,
    
    foreign key (npc_template_id) references npc_templates (id) on delete cascade,
    foreign key (item_id) references items (id) on delete cascade,
    primary key (npc_template_id, item_id)
)