CREATE TABLE IF NOT EXISTS public."employer_session"
(
    id bigserial NOT NULL,
    employer_id bigint NOT NULL,
    session_token text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT employer_session_primary_key PRIMARY KEY (id),
    CONSTRAINT employer_session_session_token_length_check CHECK (length(session_token) <= 50) NOT VALID,
    CONSTRAINT employer_session_session_token_unique UNIQUE (session_token),

    CONSTRAINT employer_session_employer_id FOREIGN KEY (employer_id)
        REFERENCES public.employer (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."employer_session";