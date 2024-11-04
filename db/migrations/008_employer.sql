CREATE TABLE IF NOT EXISTS public."employer"
(
    id bigserial NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    city_id int,
    position text NOT NULL,
    company_name_id int,
    company_description text NOT NULL,
    company_website text NOT NULL,
    path_to_profile_avatar text NOT NULL DEFAULT 'static/default_profile.png',
    contacts text,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT employer_primary_key PRIMARY KEY (id),
    CONSTRAINT employer_first_name_length_check CHECK (length(first_name) <= 50) NOT VALID,
    CONSTRAINT employer_last_name_length_check CHECK (length(last_name) <= 50) NOT VALID,
    CONSTRAINT employer_position_length_check CHECK (length(position) <= 50) NOT VALID,
    CONSTRAINT employer_company_description_length_check CHECK (length(company_description) <= 150) NOT VALID,
    CONSTRAINT employer_company_website_length_check CHECK (length(company_website) <= 50) NOT VALID,
    CONSTRAINT employer_contacts_length_check CHECK (length(contacts) <= 150) NOT VALID,
    CONSTRAINT employer_email_unique UNIQUE (email),
    CONSTRAINT employer_email_length_check CHECK (length(email) <= 50) NOT VALID,
    CONSTRAINT employer_password_hash_length_check CHECK (length(password_hash) <= 250) NOT VALID,
    CONSTRAINT employer_city_id FOREIGN KEY (city_id)
        REFERENCES public.city (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
        NOT VALID,
    CONSTRAINT employer_company_name_id FOREIGN KEY (company_name_id)
        REFERENCES public.company (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
        NOT VALID
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."employer";