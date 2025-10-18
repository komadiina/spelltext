create table quest_rewards (
    id serial primary key,
    quest_id int not null,
    reward_type varchar(32) not null,

    item_id int, -- if reward is item or consumable
    amount int, -- if reward is currency (gold/token) or xp

    constraint chk_reward_type check (
        (reward_type = 'item' and item_id is not null) or
        (reward_type in ('gold', 'token', 'xp') and amount is not null) OR 
        (reward_type = 'consumable' and item_id is not null)
    ),

    foreign key (quest_id) references quests (id)
);  

create index idx_quest_rewards_id on quest_rewards (id);
create index idx_quest_rewards_quest_id on quest_rewards (quest_id);