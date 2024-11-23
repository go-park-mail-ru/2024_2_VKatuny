CREATE TABLE IF NOT EXISTS public."user"
(
    id serial NOT NULL,
    token text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT user_token_length_check CHECK (length(token) <= 50) NOT VALID,
    CONSTRAINT user_token_unique UNIQUE (token)
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."user";