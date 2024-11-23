CREATE TABLE IF NOT EXISTS public."answer"
(
    id serial NOT NULL,
    user_id int,
    val int,
    question_id int NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT answer_primary_key PRIMARY KEY (id),
    CONSTRAINT answer_voter_id FOREIGN KEY (user_id)
        REFERENCES public.voter (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
        NOT VALID,
    CONSTRAINT answer_question_id FOREIGN KEY (question_id)
        REFERENCES public.question (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
        NOT VALID
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."answer";
