CREATE TABLE IF NOT EXISTS public."company"
(
    id serial NOT NULL,
    company_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT company_primary_key PRIMARY KEY (id),
    CONSTRAINT company_company_name_length_check CHECK (length(company_name) <= 50) NOT VALID,
    CONSTRAINT company_company_name_unique UNIQUE (company_name)
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."company";