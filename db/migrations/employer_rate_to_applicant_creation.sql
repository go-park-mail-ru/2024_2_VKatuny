CREATE TABLE IF NOT EXISTS public."employer_rate_to_applicant_creation"
(

    rate int NOT NULL,
    employer_id bigint NOT NULL,
    applicant_creation_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT employer_rate_to_applicant_creation_primary_key PRIMARY KEY (employer_id, applicant_creation_id),

    CONSTRAINT employer_rate_to_applicant_creation_employer_id FOREIGN KEY (employer_id)
        REFERENCES public.employer (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT employer_rate_to_applicant_creation_applicant_creation_id FOREIGN KEY (applicant_creation_id)
        REFERENCES public.applicant_creation (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."employer_rate_to_applicant_creation";