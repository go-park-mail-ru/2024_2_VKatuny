CREATE TABLE IF NOT EXISTS public."vacancy"
(
    id bigserial NOT NULL,
    employer_id bigint NOT NULL,
    city_id int,
    salary int NOT NULL,
    position text NOT NULL,
    vacancy_description text NOT NULL,
    work_type_id int NOT NULL DEFAULT 1,
    path_to_company_avatar text NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT vacancy_primary_key PRIMARY KEY (id),
    CONSTRAINT vacancy_position_length_check CHECK (length(position) <= 50) NOT VALID,
    CONSTRAINT vacancy_vacancy_description_length_check CHECK (length(vacancy_description) <= 2000) NOT VALID,
    CONSTRAINT vacancy_employer_id FOREIGN KEY (employer_id)
        REFERENCES public.employer (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT vacancy_work_type_id FOREIGN KEY (work_type_id)
        REFERENCES public.work_type (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET DEFAULT
        NOT VALID,
    CONSTRAINT vacancy_city_id FOREIGN KEY (city_id)
        REFERENCES public.city (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
        NOT VALID
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."vacancy";