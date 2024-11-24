-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public."position_category"
(
    id serial NOT NULL,
    category_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT position_category_primary_key PRIMARY KEY (id),
    CONSTRAINT position_category_category_name_length_check CHECK (length(category_name) <= 50) NOT VALID,
    CONSTRAINT position_category_category_name_unique UNIQUE (category_name)
);

ALTER TABLE company
    ADD fts TSVECTOR;
UPDATE company SET fts = setweight(to_tsvector('russian', "company_name"), 'A');
CREATE INDEX company_fts ON company USING GIN(fts);

CREATE OR REPLACE FUNCTION update_company_fts_function()
RETURNS TRIGGER AS $$
BEGIN
    NEW.fts = setweight(to_tsvector('russian', NEW."company_name"), 'A');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_company_fts_trigger before UPDATE OR INSERT ON company
FOR EACH ROW EXECUTE PROCEDURE update_company_fts_function();

alter table vacancy
    add position_category_id int;

alter table vacancy
    add FOREIGN KEY (position_category_id)
        REFERENCES public.position_category (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL;
insert into position_category (category_name) values ('Художник');
insert into position_category (category_name) values ('Дизайнер');
insert into position_category (category_name) values ('Музыкант');
insert into position_category (category_name) values ('Фотограф');
insert into position_category (category_name) values ('Видеограф');
insert into position_category (category_name) values ('Артист Актёр');
insert into position_category (category_name) values ('Писатель');

update vacancy  set position_category_id=2 where id=2;
update vacancy  set position_category_id=2 where id=3;
update vacancy  set position_category_id=1 where id=4;
update vacancy  set position_category_id=2 where id=5;
update vacancy  set position_category_id=2 where id=6;
update vacancy  set position_category_id=1 where id=10;




alter table cv
    add position_category_id int;

alter table cv
    add FOREIGN KEY (position_category_id)
        REFERENCES public.position_category (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
        NOT VALID;
insert into cv (applicant_id, position_rus, position_eng, job_search_status_id, cv_description, working_experience, position_category_id)
values (1, 'Художник', 'Painter', 1,  'Первоклассный художник', '10 лет опыта работы', 1);
insert into cv (applicant_id, position_rus, position_eng, job_search_status_id, cv_description, working_experience, position_category_id)
values (1, 'Музыкант', 'Musician', 1,  'Талантливый музыкант', 'Опыт работы во многих студиях записи', 3);
insert into cv (applicant_id, position_rus, position_eng, job_search_status_id, cv_description, working_experience, position_category_id)
values (3, 'Фотограф', 'Photograph', 1,  'Свадебный фотограф', 'Более 10 лет и 1000 свадеб', 4);
insert into cv (applicant_id, position_rus, position_eng, job_search_status_id, cv_description, working_experience, position_category_id)
values (2, 'Видеограф', 'Videographer', 1,  'Снимаю кино голевудского уровня', 'Более 20 лет', 5);
insert into cv (applicant_id, position_rus, position_eng, job_search_status_id, cv_description, working_experience, position_category_id)
values (4, 'Актёр', 'Actor', 1,  'Снимаюсь в кино голевудского уровня', 'Более 20 лет', 6);
insert into cv (applicant_id, position_rus, position_eng, job_search_status_id, cv_description, working_experience, position_category_id)
values (2, 'Сценарист', 'Screenwrier', 1,  'Пишу сценарии для кино голевудского уровня', 'Более 20 лет', 7);

ALTER TABLE applicant
    ADD compressed_image text;
ALTER TABLE employer
    ADD compressed_image text;
ALTER TABLE cv
    ADD compressed_image text;
ALTER TABLE vacancy
    ADD compressed_image text;
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
DROP TABLE IF EXISTS public."position_category";