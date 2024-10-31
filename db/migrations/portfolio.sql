CREATE TABLE IF NOT EXISTS public."portfolio"
(
    id bigserial NOT NULL,
    applicant_id bigint NOT NULL,
    portfolio_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT portfolio_primary_key PRIMARY KEY (id),
    CONSTRAINT portfolio_portfolio_name_length_check CHECK (length(portfolio_name) <= 50) NOT VALID,

    CONSTRAINT portfolio_applicant_id FOREIGN KEY (applicant_id)
        REFERENCES public.applicant (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."portfolio";