CREATE TABLE IF NOT EXISTS public."applicant"
(
    id serial NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    city_id text NOT NULL,
    birth_date timestamp without time zone NOT NULL,
    path_to_profile_avatar text NOT NULL DEFAULT 'static/default_profile.png',
    contacts text,
    education text,
    email text NOT NULL,
    password_hash text NOT NULL,
    description text,
    CONSTRAINT applicant_primary_key PRIMARY KEY (id),
    CONSTRAINT applicant_first_name_length_check CHECK (length(first_name) <= 50) NOT VALID,
    CONSTRAINT applicant_last_name_length_check CHECK (length(last_name) <= 50) NOT VALID,
    CONSTRAINT applicant_contacts_length_check CHECK (length(contacts) <= 150) NOT VALID,
    CONSTRAINT applicant_education_length_check CHECK (length(education) <= 150) NOT VALID,
    CONSTRAINT applicant_email_key UNIQUE (email),
    CONSTRAINT applicant_email_length_check CHECK (length(email) <= 50) NOT VALID,
    CONSTRAINT applicant_password_hash_length_check CHECK (length(password_hash) <= 250) NOT VALID,
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."applicant";