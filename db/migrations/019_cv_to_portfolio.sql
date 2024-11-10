CREATE TABLE IF NOT EXISTS public."cv_to_portfolio"
(
    cv_id bigint NOT NULL,
    portfolio_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT cv_to_portfolio_primary_key PRIMARY KEY (cv_id, portfolio_id),
    CONSTRAINT cv_to_portfolio_cv_id FOREIGN KEY (cv_id)
        REFERENCES public.cv (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT cv_to_portfolio_portfolio_id FOREIGN KEY (portfolio_id)
        REFERENCES public.portfolio (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
    
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."cv_to_portfolio";