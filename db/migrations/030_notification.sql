-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public."employer_notification"
(
    id serial NOT NULL,
    notification_text text NOT NULL,
    employer_id bigint NOT NULL,
    vacancy_id bigint NOT NULL,
    applicant_id bigint NOT NULL,
    is_read boolean DEFAULT false NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT employer_notification_primary_key PRIMARY KEY (id),
    CONSTRAINT employer_notification_notification_text_length_check CHECK (length(notification_text) <= 200) NOT VALID
);
insert into employer_notification (notification_text, employer_id, vacancy_id, applicant_id) values ('На вашу вакансию "Дизайнер витрины" откликнулся Иван Иванов', 1, 1, 1);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
DROP TABLE IF EXISTS public."employer_notification";