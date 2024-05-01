CREATE TABLE IF NOT EXISTS tags_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_ character varying(50),
    slug_ character varying(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS articles_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title_ character varying(255),
    subtitle_ character varying(255),
    slug_ character varying(50) UNIQUE NOT NULL,
    body_ text,
    created_at_ timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at_ timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS article_tags_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    article_id_ integer,
    tag_id_ integer references tags_(id_) NOT NULL
);

-- select a.
--
-- select i.title, i.image, i.type, i.id, r.title, r.body, r.recipe_id from
--   recipe_ingredients ri
--   inner join
--   ingredients i
--   on i.id = ri.ingredient_id
--   inner join
--   recipes r
--   on r.recipe_id = ri.recipe_id
--   where r.recipe_id = 1;
