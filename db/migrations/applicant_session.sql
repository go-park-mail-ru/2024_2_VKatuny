CREATE TABLE IF NOT EXISTS public."applicant_session"
(
    id bigserial NOT NULL,
    applicant_id bigint NOT NULL,
    session_token text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT applicant_session_primary_key PRIMARY KEY (id),
    CONSTRAINT applicant_session_session_token_length_check CHECK (length(session_token) <= 50) NOT VALID,
    CONSTRAINT applicant_session_session_token_unique UNIQUE (session_token),

    CONSTRAINT applicant_session_applicant_id FOREIGN KEY (applicant_id)
        REFERENCES public.applicant (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."applicant_session";