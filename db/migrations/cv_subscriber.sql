CREATE TABLE IF NOT EXISTS public."cv_subscriber"
(
    cv_id bigint NOT NULL,
    employer_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT cv_subscriber_primary_key PRIMARY KEY (cv_id, employer_id),

    CONSTRAINT cv_subscriber_cv_id FOREIGN KEY (cv_id)
        REFERENCES public.cv (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT cv_subscriber_employer_id FOREIGN KEY (employer_id)
        REFERENCES public.employer (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."cv_subscriber";