CREATE TABLE IF NOT EXISTS public."cv"
(
    id bigserial NOT NULL,
    applicant_id bigint NOT NULL,
    position_rus text NOT NULL,
    position_eng text NOT NULL,
    job_search_status_id int NOT NULL DEFAULT 1,
    working_experience text NOT NULL,
    path_to_profile_avatar text NOT NULL DEFAULT 'static/default_profile.png',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT cv_primary_key PRIMARY KEY (id),
    CONSTRAINT cv_position_rus_length_check CHECK (length(position_rus) <= 50) NOT VALID,
    CONSTRAINT cv_position_eng_length_check CHECK (length(position_eng) <= 50) NOT VALID,
    CONSTRAINT cv_working_experience_length_check CHECK (length(working_experience) <= 1000) NOT VALID,
    CONSTRAINT cv_applicant_id FOREIGN KEY (applicant_id)
        REFERENCES public.applicant (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT cv_job_search_status_id FOREIGN KEY (job_search_status_id)
        REFERENCES public.job_search_status (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET DEFAULT
        NOT VALID
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."cv";