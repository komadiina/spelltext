create table npcs (
    id serial primary key,
    prefix varchar(64),
    suffix varchar(64),
    template_id int not null,

    health_multiplier float not null default 1.0,
    damage_multiplier float not null default 1.0,

    foreign key (template_id) references npc_templates (id)
);

create index idx_npcs_id on npcs (id);

-- alter table npcs replica identity full;