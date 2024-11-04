CREATE TABLE IF NOT EXISTS public."job_search_status"
(
    id serial NOT NULL,
    job_search_status_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT job_search_status_primary_key PRIMARY KEY (id),
    CONSTRAINT job_search_status_job_search_status_name_length_check CHECK (length(job_search_status_name) <= 50) NOT VALID,
    CONSTRAINT job_search_status_job_search_status_name_unique UNIQUE (job_search_status_name)
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."job_search_status";