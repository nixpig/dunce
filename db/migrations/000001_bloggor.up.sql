CREATE TABLE IF NOT EXISTS type_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_ character varying(255) NOT NULL,
    template_ character varying(255) NOT NULL,
    slug_ character varying(255) NOT NULL
);

INSERT INTO type_ (name_, template_, slug_) VALUES ('post', 'post', 'post');

CREATE TABLE IF NOT EXISTS role_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_ character varying(100) NOT NULL
);

INSERT INTO role_ (name_) VALUES ('admin'), ('author'), ('reader');

CREATE TABLE IF NOT EXISTS tag_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_ character varying(100) NOT NULL,
    slug_ character varying(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username_ character varying(100) NOT NULL,
    email_ character varying(100) NOT NULL,
    password_ character varying(255) NOT NULL,
    link_ character varying(255) NOT NULL,
    role_ integer references role_(id_) NOT NULL
);

CREATE TABLE IF NOT EXISTS site_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_ character varying(100),
    description_ character varying(255),
    url_ character varying(255),
    owner_ integer references user_(id_) NOT NULL
);

CREATE TABLE IF NOT EXISTS article_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title_ character varying(255),
    subtitle_ character varying(255),
    slug_ character varying(255) NOT NULL,
    body_ text,
    created_at_ timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at_ timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    type_id_ integer references type_(id_) NOT NULL,
    user_id_ integer references user_(id_) NOT NULL,
    tag_ids_ character varying(255)
);

