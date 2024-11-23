CREATE TABLE IF NOT EXISTS public."question"
(
    id serial NOT NULL,
    question_text text,
    type_id int NOT NULL,
    position int NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT question_primary_key PRIMARY KEY (id),
    CONSTRAINT question_type_id FOREIGN KEY (type_id)
        REFERENCES public.question_type (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
        NOT VALID
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."question";
