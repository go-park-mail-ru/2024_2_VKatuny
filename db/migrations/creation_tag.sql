CREATE TABLE IF NOT EXISTS public."creation_tag"
(
    id serial NOT NULL,
    creation_tag_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT creation_tag_primary_key PRIMARY KEY (id),
    CONSTRAINT creation_tag_creation_tag_name_length_check CHECK (length(creation_tag_name) <= 50) NOT VALID,
    CONSTRAINT creation_tag_creation_tag_name_unique UNIQUE (creation_tag_name),
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."creation_tag";