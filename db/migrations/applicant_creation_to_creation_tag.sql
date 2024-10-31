CREATE TABLE IF NOT EXISTS public."applicant_creation_to_creation_tag"
(
    creation_tag_id bigint NOT NULL,
    applicant_creation_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT applicant_creation_to_creation_tag_primary_key PRIMARY KEY (creation_tag_id, applicant_creation_id),

    CONSTRAINT applicant_creation_to_creation_tag_creation_tag_id FOREIGN KEY (creation_tag_id)
        REFERENCES public.creation_tag (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET 1
        NOT VALID,
    CONSTRAINT applicant_creation_to_creation_tag_applicant_creation_id FOREIGN KEY (applicant_creation_id)
        REFERENCES public.applicant_creation (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."applicant_creation_to_creation_tag";