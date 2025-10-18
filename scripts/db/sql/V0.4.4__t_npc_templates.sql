create table npc_templates (
    id serial primary key,
    name varchar(64) not null,

    min_level int not null,
    max_level int not null,

    health_points int not null,
    base_damage int not null,
    base_xp_reward int not null,

    drop_item_id int,
    
    constraint xp_reward_positive check (base_xp_reward > 0),
    constraint valid_level check (min_level <= max_level)
);

create index idx_npc_templates_id on npc_templates (id);
create index idx_npc_templates_drop_item_id on npc_templates (drop_item_id);
create index idx_levels on npc_templates (min_level, max_level);