create table quest_prerequisites (
    base_quest_id int not null,
    prerequisite_quest_id int not null,
    
    foreign key (base_quest_id) references quests (id),
    foreign key (prerequisite_quest_id) references quests (id)
);