CREATE TABLE IF NOT EXISTS types_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_ character varying(255) UNIQUE NOT NULL,
    template_ character varying(255) NOT NULL,
    slug_ character varying(255) UNIQUE NOT NULL
);

INSERT INTO types_ (name_, template_, slug_) VALUES ('post', 'pages/public/post', 'posts'), ('page', 'pages/public/page', 'page');

CREATE TYPE role_ AS ENUM ('admin', 'author', 'reader');

CREATE TABLE IF NOT EXISTS tags_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_ character varying(50) UNIQUE NOT NULL,
    slug_ character varying(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username_ character varying(100) UNIQUE NOT NULL,
    email_ character varying(100) UNIQUE NOT NULL,
    password_ character varying(255) NOT NULL,
    link_ character varying(255) NOT NULL,
    role_ role_ NOT NULL
);

insert into users_ (username_, email_, password_, link_, role_) values ('admin', 'admin@example.org', 'p4ssw0rd', '', 'admin');

CREATE TABLE IF NOT EXISTS site_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_ character varying(100),
    description_ character varying(255),
    url_ character varying(255),
    owner_ integer references users_(id_) NOT NULL
);

CREATE TABLE IF NOT EXISTS articles_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title_ character varying(255),
    subtitle_ character varying(255),
    slug_ character varying(255) UNIQUE NOT NULL,
    body_ text,
    created_at_ timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at_ timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    type_id_ integer references types_(id_) NOT NULL,
    user_id_ integer references users_(id_) NOT NULL,
    tag_ids_ character varying(255)
);
