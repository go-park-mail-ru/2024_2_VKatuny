CREATE TABLE IF NOT EXISTS public."vacancy_subscriber"
(
    cv_id bigint NOT NULL,
    applicant_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT vacancy_subscriber_primary_key PRIMARY KEY (cv_id, applicant_id),
    CONSTRAINT vacancy_subscriber_cv_id FOREIGN KEY (cv_id)
        REFERENCES public.cv (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT vacancy_subscriber_applicant_id FOREIGN KEY (applicant_id)
        REFERENCES public.employer (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
    
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."vacancy_subscriber";