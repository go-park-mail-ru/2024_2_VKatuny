-- Write your migrate up statements here
insert into city (id, city_name) values (1, 'Москва');
insert into company (id, company_name) values (1, 'Яндекс');
insert into creation_tag (id, creation_tag_name) values (1, 'Скульптура');
insert into creation_type (id, creation_type_name) values (1, 'svg');
insert into job_search_status (id, job_search_status_name) values (1, 'Активно ищу работу');
insert into applicant (id, first_name , last_name , city_id, birth_date, path_to_profile_avatar, contacts, education, email, password_hash)
values (1, 'Иван', 'Иванов', 1, '12-12-2001', 'IvanovIvan.svg', 'tg - @IvanovIvan', 'МАРХИ', 'ivanovivan@mail.ru', '23e16c2f9af678bd3485cb9ccf8b6e8e');
insert into employer (id, first_name , last_name , city_id, position, company_name_id, company_description, company_website, path_to_profile_avatar, contacts,	email, password_hash)
values (1, 'Петр', 'Петров', 1, 'Помощник мэра', 1, 'Мэрия Москвы', 'https://www.mos.ru/', 'PetrPetrov.svg', 'tg - @PetrPetrov', 'petrovpetr@mail.ru',  'a95e46339b08617435e0cf390bcce5ff');
insert into work_type (id, work_type_name) values (1, 'Полная занятость');
insert into cv (id, applicant_id, position_rus, position_eng, job_search_status_id, cv_description, working_experience, path_to_profile_avatar)
values (1, 1, 'Скульптор', 'Sculptor', 1,  'Я усердный и целеустремленный', 'Не было опыта работы', 'IvanovIvan.svg');
insert into cv_subscriber (employer_id , cv_id) values (1, 1);
insert into applicant_session (id, applicant_id, session_token) values (1, 1, '1ahsdfybegtorhjertoldbtsdjgxsdfkg');
insert into employer_session (id , employer_id , session_token) values (1, 1, '2heysdfyuilsorhjertuhebtsdjxsdfkg');
insert into portfolio (id, applicant_id, portfolio_name) values (1, 1, 'Моё скульптурное портфолио');
insert into vacancy (id, employer_id , salary, position, vacancy_description, work_type_id, path_to_company_avatar,city_id)
values (1, 1, 90000, 'Скульптор', 'Требуется скульптор без опыта работы', 1, 'logoBMSTU.svg', 1);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (1, 1);
insert into applicant_creation (id , applicant_id , applicant_creation_name, path_to_creation, creation_type_id)
values (1, 1, 'Скульптура сократа', 'Socrates.svg', 1);
insert into applicant_rate_to_applicant_creation (rate, applicant_id, applicant_creation_id) values (4, 1, 1);
insert into cv_to_portfolio (cv_id, portfolio_id) values (1, 1);
insert into vacancy_to_creation_tag (vacancy_id , creation_tag_id) values (1, 1);
insert into employer_rate_to_applicant_creation (rate, employer_id, applicant_creation_id) values (5, 1, 1);
insert into cv_to_creation_tag (cv_id , creation_tag_id) values (1, 1);
insert into applicant_creation_to_creation_tag (applicant_creation_id , creation_tag_id) values (1, 1);
insert into applicant_creation_to_portfolio (applicant_creation_id , portfolio_id)  values (1, 1);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
delete from city where city_name='Москва';
delete from company where company_name='Яндекс';
delete from creation_tag where creation_tag_name='Скульптура';
delete from creation_type where creation_type_name='svg';
delete from job_search_status where job_search_status_name='Активно ищу работу';
delete from applicant where email='ivanovivan@mail.ru';
delete from employer where email='petrovpetr@mail.ru';
delete from work_type where work_type_name='Полная занятость';
delete from cv where id=1;
delete from cv_subscriber where employer_id=1 and cv_id=1;
delete from applicant_session where session_token='1ahsdfybegtorhjertoldbtsdjgxsdfkg';
delete from employer_session where session_token='2heysdfyuilsorhjertuhebtsdjxsdfkg';
delete from portfolio where id=1;
delete from vacancy where id=1;
delete from vacancy_subscriber where applicant_id=1 and vacancy_id=1;
delete from applicant_creation where id=1;
delete from applicant_rate_to_applicant_creation where applicant_id=1 and applicant_creation_id=1;
delete from cv_to_portfolio where cv_id=1 and portfolio_id=1;
delete from vacancy_to_creation_tag where vacancy_id=1 and creation_tag_id=1;
delete from employer_rate_to_applicant_creation where employer_id=1 and applicant_creation_id=1;
delete from cv_to_creation_tag where cv_id=1 and creation_tag_id=1;
delete from cv_to_creation_tag where applicant_creation_id=1 and creation_tag_id=1;
delete from cv_to_creation_tag where applicant_creation_id=1 and portfolio_id=1;