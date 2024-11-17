-- Write your migrate up statements here
alter table vacancy
    add fts TSVECTOR
;
CREATE INDEX vacancy_fts ON vacancy  USING GIN(fts);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

drop INDEX vacancy_fts;