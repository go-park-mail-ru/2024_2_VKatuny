CREATE TABLE IF NOT EXISTS public."vacancy_subscriber"
(
    vacancy_id bigint NOT NULL,
    applicant_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT vacancy_subscriber_primary_key PRIMARY KEY (vacancy_id, applicant_id),
    CONSTRAINT vacancy_subscriber_vacancy_id FOREIGN KEY (vacancy_id)
        REFERENCES public.vacancy (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT vacancy_subscriber_applicant_id FOREIGN KEY (applicant_id)
        REFERENCES public.applicant (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
    
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."vacancy_subscriber";