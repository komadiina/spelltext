create table quest_objective (
    quest_id int not null,
    npc_id int not null,
    foreign key (quest_id) references quests (id),
    foreign key (npc_id) references npcs (id)
);  

-- alter table quest_objective replica identity full;