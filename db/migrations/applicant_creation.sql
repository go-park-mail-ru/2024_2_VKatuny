CREATE TABLE IF NOT EXISTS public."applicant_creation"
(
    id serial NOT NULL,
    applicant_id bigint NOT NULL,
    applicant_creation_name text NOT NULL,
    path_to_creation  text NOT NULL,
    creation_type_id int NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT applicant_creation_primary_key PRIMARY KEY (id),
    CONSTRAINT applicant_creation_applicant_creation_name_length_check CHECK (length(applicant_creation_name) <= 50) NOT VALID,
    CONSTRAINT applicant_creation_path_to_creation_unique UNIQUE (path_to_creation),
    
    CONSTRAINT applicant_creation_applicant_id FOREIGN KEY (applicant_id)
        REFERENCES public.applicant (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT applicant_creation_creation_type_id FOREIGN KEY (creation_type_id)
        REFERENCES public.creation_type (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET 1
        NOT VALID,
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."applicant_creation";