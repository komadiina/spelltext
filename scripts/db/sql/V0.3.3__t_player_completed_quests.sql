create table character_completed_quests (
    character_id int not null,
    quest_id int not null,
    
    primary key (character_id, quest_id),
    foreign key (character_id) references characters (character_id) on delete cascade,
    foreign key (quest_id) references quests (id)
);

-- alter table character_completed_quests replica identity full;