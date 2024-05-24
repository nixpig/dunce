create table if not exists site_ (
    id_ integer primary key generated always as identity,
    key_ character varying(50) unique not null,
    value_ character varying(255) not null,
    meta_ character varying(255) not null
);
