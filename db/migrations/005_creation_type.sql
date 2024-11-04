CREATE TABLE IF NOT EXISTS public."creation_type"
(
    id serial NOT NULL,
    creation_type_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT creation_type_primary_key PRIMARY KEY (id),
    CONSTRAINT creation_type_creation_type_name_length_check CHECK (length(creation_type_name) <= 50) NOT VALID,
    CONSTRAINT creation_type_creation_type_name_unique UNIQUE (creation_type_name)
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."creation_type";