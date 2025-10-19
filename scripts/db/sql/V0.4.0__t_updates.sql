create table updates (
    id serial primary key,
    title varchar(64) not null,
    version varchar(16) not null,
    description text not null
);

create index idx_version on updates (version);
alter table updates replica identity full;