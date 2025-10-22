create table npc_abilities (
    npc_id int not null,
    ability_id int not null,

    foreign key (npc_id) references npcs (id),
    foreign key (ability_id) references abilities (id),
    primary key (npc_id, ability_id)
);

create index idx_npc_abilities_npc_id on npc_abilities (npc_id);
create index idx_npc_abilities_ability_id on npc_abilities (ability_id);