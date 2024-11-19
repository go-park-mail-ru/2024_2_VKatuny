-- Write your migrate up statements here

ALTER TABLE vacancy
    ADD fts TSVECTOR
;
ALTER TABLE cv
    ADD fts TSVECTOR
;

UPDATE cv SET fts = setweight(to_tsvector('russian', "position_rus"), 'A')|| setweight(to_tsvector('english', "position_eng"), 'A') || setweight(to_tsvector('russian', cv_description), 'B') || setweight(to_tsvector('russian', working_experience), 'B');
UPDATE vacancy SET fts = setweight(to_tsvector('russian', "position"), 'A') || setweight(to_tsvector('russian', vacancy_description), 'B');

CREATE INDEX cv_fts ON cv  USING GIN(fts);
CREATE INDEX vacancy_fts ON vacancy  USING GIN(fts);

CREATE OR REPLACE FUNCTION update_cv_fts_function()
RETURNS TRIGGER AS $$
BEGIN
    NEW.fts = setweight(to_tsvector('russian', NEW."position_rus"), 'A')|| setweight(to_tsvector('english', NEW."position_eng"), 'A') || setweight(to_tsvector('russian', NEW.cv_description), 'B') || setweight(to_tsvector('russian', NEW.working_experience), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE FUNCTION update_vacancy_fts_function()
RETURNS TRIGGER AS $$
BEGIN
    NEW.fts = setweight(to_tsvector('russian', NEW."position"), 'A') || setweight(to_tsvector('russian', NEW.vacancy_description), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_cv_fts_trigger before UPDATE OR INSERT ON cv
FOR EACH ROW EXECUTE PROCEDURE update_cv_fts_function();
CREATE TRIGGER update_vacancy_fts_trigger before UPDATE OR INSERT ON vacancy
FOR EACH ROW EXECUTE PROCEDURE update_vacancy_fts_function();
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

DROP INDEX vacancy_fts;
DROP INDEX cv_fts;
DROP TRIGGER update_cv_fts_trigger;
DROP TRIGGER update_vacancy_fts_trigger;
DROP FUNCTION update_cv_fts_function;
DROP FUNCTION update_vacancy_fts_function;