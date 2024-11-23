-- Write your migrate up statements here
insert into question_type (question_type_name) values ('CSAT');
insert into question (question_text, type_id, position) values ('Насколько вам понравилось пользоваться нашим сервисом?', 1, 1);
insert into question (question_text, type_id, position) values ('Насколько частовы пользуетесь нашим сервисом?', 1, 2);
insert into question (question_text, type_id, position) values ('Насколько полезным кажется наш сервис для решения ваших задач?', 1, 3);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
delete from question_type where question_type_name='CSAT';
delete from question where id=1;
delete from question where id=2;
delete from question where id=3;