CREATE TABLE IF NOT EXISTS public."vacancy_to_creation_tag"
(
    creation_tag_id bigint NOT NULL,
    vacancy_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT vacancy_to_creation_tag_primary_key PRIMARY KEY (creation_tag_id, vacancy_id),

    CONSTRAINT vacancy_to_creation_tag_creation_tag_id FOREIGN KEY (creation_tag_id)
        REFERENCES public.creation_tag (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET 1
        NOT VALID,
    CONSTRAINT vacancy_to_creation_tag_vacancy_id FOREIGN KEY (vacancy_id)
        REFERENCES public.vacancy (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."vacancy_to_creation_tag";