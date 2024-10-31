CREATE TABLE IF NOT EXISTS public."applicant_creation_to_portfolio"
(
    portfolio_id bigint NOT NULL,
    applicant_creation_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT applicant_creation_to_portfolio_primary_key PRIMARY KEY (portfolio_id, applicant_creation_id),

    CONSTRAINT applicant_creation_to_portfolio_portfolio_id FOREIGN KEY (portfolio_id)
        REFERENCES public.portfolio (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT applicant_creation_to_portfolio_applicant_creation_id FOREIGN KEY (applicant_creation_id)
        REFERENCES public.applicant_creation (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."applicant_creation_to_portfolio";