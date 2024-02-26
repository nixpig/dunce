CREATE TABLE IF NOT EXISTS sessions_ (
    id_ integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    token_ text,
    expires_at_ bigint,
    issued_at_ bigint,
    user_id_ integer references users_(id_) NOT NULL
);
