-- Write your migrate up statements here
insert into city (city_name) values ('Москва');
insert into city (city_name) values ('Санкт-Петербург');
insert into city (city_name) values ('Казань');
insert into company (company_name) values ('Мэрия Москвы');
insert into company (company_name) values ('Архитект Бутик');
insert into company (company_name) values ('Группа Пик');
insert into company (company_name) values ('Teikaboom');
insert into creation_tag (creation_tag_name) values ('Скульптура');
insert into creation_type (creation_type_name) values ('svg');
insert into job_search_status (job_search_status_name) values ('Активно ищу работу');
insert into applicant (first_name , last_name , city_id, birth_date, path_to_profile_avatar, contacts, education, email, password_hash)
values ('Иван', 'Иванов', 1, '12-12-2001', '/media/Uncompressed/1ahsdfybegtorhlodjtldbtsdjgxsdfkg.JPG', 'tg - @IvanovIvan', 'МАРХИ', 'applicant@mail.ru', '$2a$10$mxG.iijgPJyg3RXdCdDyT.Nrah32oBs5JfaIoum4ITx.PMF.oNV1a'); --pass1234

insert into applicant (first_name , last_name , city_id, birth_date, email, password_hash)
values ('Дмитрий', 'Петров', 1, '12-10-1989', 'applicant1@mail.ru', '$2a$10$mxG.iijgPJyg3RXdCdDyT.Nrah32oBs5JfaIoum4ITx.PMF.oNV1a'); --pass1234
insert into applicant (first_name , last_name , city_id, birth_date, email, password_hash)
values ('Максим', 'Тихомиров', 1, '09-02-2003', 'applicant2@mail.ru', '$2a$10$mxG.iijgPJyg3RXdCdDyT.Nrah32oBs5JfaIoum4ITx.PMF.oNV1a'); --pass1234
insert into applicant (first_name , last_name , city_id, birth_date, email, password_hash)
values ('Сергей', 'Волков', 1, '12-08-2001', 'applicant3@mail.ru', '$2a$10$mxG.iijgPJyg3RXdCdDyT.Nrah32oBs5JfaIoum4ITx.PMF.oNV1a'); --pass1234
insert into applicant (first_name , last_name , city_id, birth_date, email, password_hash)
values ('Владимир', 'Ершов', 1, '12-08-2001', 'applicant4@mail.ru', '$2a$10$mxG.iijgPJyg3RXdCdDyT.Nrah32oBs5JfaIoum4ITx.PMF.oNV1a'); --pass1234

insert into employer (first_name , last_name , city_id, position, company_name_id, company_description, company_website, path_to_profile_avatar, contacts,	email, password_hash)
values ('Петр', 'Петров', 1, 'Помощник мэра', 1, 'Мэрия Москвы', 'https://www.mos.ru/', '/media/Uncompressed/1uytdfybegtorhlodjtldbtsdjioldfkg.JPG', 'tg - @PetrPetrov', 'employer@mail.ru',  '$2a$10$UJdgr8sjQPsa1IpS7pLHBu3VgsO4W/SPjGjVBI2aw1WdYcx63IAEK'); --pass1234

insert into employer (first_name , last_name , city_id, position, company_name_id, company_description, company_website,	email, password_hash)
values ('Станислав', 'Шубин', 2, 'Ландшафтный дизайнер', 2, 'Архитект Бутик', 'https://architect.boutique/', 'employer2@mail.ru',  '$2a$10$UJdgr8sjQPsa1IpS7pLHBu3VgsO4W/SPjGjVBI2aw1WdYcx63IAEK'); --pass1234
insert into employer (first_name , last_name , city_id, position, company_name_id, company_description, company_website,	email, password_hash)
values ('Илья', 'Миронов', 3, 'Дизайнер интерьера', 3, 'Группа Пик', 'https://www.pik.ru/', 'employer3@mail.ru',  '$2a$10$UJdgr8sjQPsa1IpS7pLHBu3VgsO4W/SPjGjVBI2aw1WdYcx63IAEK'); --pass1234
insert into employer (first_name , last_name , city_id, position, company_name_id, company_description, company_website,	email, password_hash)
values ('Семен', 'Стрельцов', 1, 'Главный дизайнер', 4, 'Teikaboom', 'https://teikaboom.ru/', 'employer4@mail.ru',  '$2a$10$UJdgr8sjQPsa1IpS7pLHBu3VgsO4W/SPjGjVBI2aw1WdYcx63IAEK'); --pass1234

insert into work_type (work_type_name) values ('Полная занятость');
insert into work_type (work_type_name) values ('Разовая работа');
insert into cv (applicant_id, position_rus, position_eng, job_search_status_id, cv_description, working_experience, path_to_profile_avatar)
values (1, 'Скульптор', 'Sculptor', 1,  'Я усердный и целеустремленный', 'Не было опыта работы', '/media/Uncompressed/1ahsdfybegtorhlodjtldbtsdjgxsdfkg.JPG');
insert into cv_subscriber (employer_id , cv_id) values (1, 1);
insert into applicant_session ( applicant_id, session_token) values (1, '1ahsdfybegtorhjertoldbtsdjgxsdfkg');
insert into employer_session (employer_id , session_token) values (1, '2heysdfyuilsorhjertuhebtsdjxsdfkg');
insert into portfolio (applicant_id, portfolio_name) values (1, 'Моё скульптурное портфолио');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id)
values (1, 90000, 'Скульптор', 'Требуется скульптор без опыта работы', 1, 1);


insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (1, 170000, 'Дизайнер витрины', 'Требуется оформить главный стенд нового офиса', 2, 2, 'media/Uncompressed/1ahsdfhtrgtorhjertoldbtsdjgxsdfkg.PNG');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (3, 80000, 'Младший дизайнер интерьера', 'Требуется специалист в области дизайна интерьера в дружный коллектив нашего бюро', 1, 2, 'media/Uncompressed/1ahsdfhtrgtaahjertoldbtsdjgxsdfkg.PNG');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (3, 100000, 'Художник оформитель', 'Требуется опытный специалист для оформления дизайнерской мебели', 1, 3, 'media/Uncompressed/1ahsdfhtrgtaahjertoldbtsdjgxsdfkg.PNG');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (2, 210000, 'Младший ландшафтный дизайнер', 'Нанимаем опытного специалиста для оформления нескольких парковых зон', 2, 1, 'media/Uncompressed/1ahsdfhtrgtorhjruooldbtsdjgxsdfkg.JPG');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (2, 120000, 'Младший ландшафтный дизайнер', 'Требуется специалист по живой изгороди', 1, 3, 'media/Uncompressed/1ahsdfhtrgtorhjruooldbtsdjgxsdfkg.JPG');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (2, 220000, 'Куратор', 'Требуется куратор в нашу программу повышения квалификации сотрудников', 1, 2, 'media/Uncompressed/1ahsdfhtrgtorhjruooldbtsdjgxsdfkg.JPG');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (1, 320000, 'Графитист', 'Нанимаем профессионального художника-гафитиста для временного графити приуроченному к празднику на высотном здании', 2, 2, 'media/Uncompressed/1ahsdfhtrgtorhjertoldbtsdjgxsdfkg.PNG');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (4, 110000, 'Хореограф', 'Ищем профессионального хареографа со стажем 15 лет', 1, 3, 'media/Uncompressed/1ahsdfhtrgtorhjruaaldbtsdjgxsdfkg.JPG');
insert into vacancy (employer_id , salary, position, vacancy_description, work_type_id, city_id, path_to_company_avatar)
values (4, 210000, 'Художник декоратор', 'Требуется декоратор для постановок', 1, 1, 'media/Uncompressed/1ahsdfhtrgtorhjruaaldbtsdjgxsdfkg.JPG');


insert into vacancy_subscriber (applicant_id , vacancy_id) values (1, 1);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (2, 1);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (3, 1);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (1, 2);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (2, 2);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (3, 2);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (4, 2);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (5, 2);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (3, 3);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (2, 3);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (4, 4);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (5, 4);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (1, 5);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (2, 6);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (4, 7);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (5, 8);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (1, 8);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (2, 9);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (4, 9);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (2, 10);
insert into vacancy_subscriber (applicant_id , vacancy_id) values (1, 10);
insert into applicant_creation (applicant_id , applicant_creation_name, path_to_creation, creation_type_id)
values (1, 'Скульптура сократа', 'Socrates.svg', 1);
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
