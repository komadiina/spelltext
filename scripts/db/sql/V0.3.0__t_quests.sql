create table quests (
    id serial primary key,
    name varchar(64) not null,
    description text not null,
    level_requirement int not null default 1
);

create index idx_quests_id on quests (id);

-- alter table quests replica identity full;