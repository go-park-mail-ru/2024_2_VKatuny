CREATE TABLE IF NOT EXISTS public."favorite_vacancy"
(
    applicant_id bigint NOT NULL,
    vacancy_id bigint NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT favorite_vacancy_primary_key PRIMARY KEY (applicant_id, vacancy_id),
    CONSTRAINT favorite_vacancy_applicant_id FOREIGN KEY (applicant_id)
        REFERENCES public.applicant (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT favorite_vacancy_vacancy_id FOREIGN KEY (vacancy_id)
        REFERENCES public.vacancy (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
    
);
CREATE TRIGGER update_applicant_favorite_vacancy before UPDATE ON favorite_vacancy
FOR EACH ROW EXECUTE PROCEDURE update_time();
insert into favorite_vacancy (applicant_id , vacancy_id) values (1, 1);
insert into favorite_vacancy (applicant_id , vacancy_id) values (2, 2);
insert into favorite_vacancy (applicant_id , vacancy_id) values (3, 3);
---- create above / drop below ----

DROP TABLE IF EXISTS public."favorite_vacancy";
DROP TRIGGER update_applicant_favorite_vacancy;