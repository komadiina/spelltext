create table gamba_chest_contents (
    gamba_chest_id int not null,
    item_id int not null,
    
    foreign key (gamba_chest_id) references gamba_chests (id),
    foreign key (item_id) references items (id)
)

create index idx_gamba_chest_contents_item_id on gamba_chest_contents (item_id);