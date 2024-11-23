CREATE TABLE IF NOT EXISTS public."question_type"
(
    id serial NOT NULL,
    question_type_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT question_type_primary_key PRIMARY KEY (id),
    CONSTRAINT question_type_question_type_name_length_check CHECK (length(question_type_name) <= 50) NOT VALID,
    CONSTRAINT question_type_question_type_name_unique UNIQUE (question_type_name)
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."question_type";