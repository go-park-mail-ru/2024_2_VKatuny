CREATE TABLE IF NOT EXISTS public."cv_to_creation_tag"
(
    creation_tag_id bigint NOT NULL,
    cv_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT cv_to_creation_tag_primary_key PRIMARY KEY (creation_tag_id, cv_id),
    CONSTRAINT cv_to_creation_tag_creation_tag_id FOREIGN KEY (creation_tag_id)
        REFERENCES public.creation_tag (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT cv_to_creation_tag_cv_id FOREIGN KEY (cv_id)
        REFERENCES public.cv (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
    
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."cv_to_creation_tag";