create table npc_templates (
    id serial primary key,
    name varchar(64) not null,
    description varchar(255) not null default 'empty description',

    min_level int not null,
    max_level int not null,

    health_points int not null,
    base_damage int not null,
    base_xp_reward int not null,
    gold_reward int not null default 10,
    
    constraint xp_reward_positive check (base_xp_reward > 0),
    constraint valid_level check (min_level <= max_level)
);

create index idx_npc_templates_id on npc_templates (id);
create index idx_levels on npc_templates (min_level, max_level);

-- alter table npc_templates replica identity full;