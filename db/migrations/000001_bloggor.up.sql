create table if not exists tags_ (
    id_ integer primary key generated always as identity,
    name_ character varying(50),
    slug_ character varying(50) unique not null
);

create table if not exists articles_ (
    id_ integer primary key generated always as identity,
    title_ character varying(255),
    subtitle_ character varying(255),
    slug_ character varying(50) unique not null,
    body_ text,
    created_at_ timestamp without time zone default current_timestamp not null,
    updated_at_ timestamp without time zone default current_timestamp not null
);

create table if not exists article_tags_ (
    id_ integer primary key generated always as identity,
    article_id_ integer,
    tag_id_ integer references tags_(id_) not null
);

create table if not exists users_ (
    id_ integer primary key generated always as identity,
    name_ character varying(255) unique not null,
    email_ character varying(255) unique not null,
    password_ character varying(60) not null,
    created_at_ timestamp without time zone default current_timestamp not null
);
