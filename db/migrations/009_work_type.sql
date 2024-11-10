CREATE TABLE IF NOT EXISTS public."work_type"
(
    id serial NOT NULL,
    work_type_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT work_type_primary_key PRIMARY KEY (id),
    CONSTRAINT work_type_work_type_name_length_check CHECK (length(work_type_name) <= 50) NOT VALID,
    CONSTRAINT work_type_work_type_name_unique UNIQUE (work_type_name)
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."work_type";