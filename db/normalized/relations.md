# Описание Базы Данных

## Первая нормальная форма
База данных находится в первой нормальной форме, так как все ее атрибуты являются простыми (строчки не содержат json списиски или любые другие перечисления), не содержит повторяющихся строк, что обеспечивает независимость данных от номера строки.

## Второая нормальная форма
База данных находится во второй нормальной форме, так как находится в первой нормальной форме и каждый неключевой атребут неприводимо зависит от первичного ключа (нельзя выявить зависимость неключевого поля от другого неключевого поля).

## Третья нормальная форма
База данных находится в третьей нормальной форме, так как находится во второй нормальной форме и каждый неключевой атрибут нетранзитивно зависит от первичного ключа (Все атрибуты, значение которых часто повторяется в разных строчках или вовсе ограничено некоторым множеством значений вынесены в отдельные таблицы, а первичные ключи синтетические).

## Нормальная форма Бойса-Кодда
База данных находится в нормальной форме Бойса-Кодда, так как находится в третьей нормальной форме и каждая нетривиальная и неприводимая слева функциональная зависимость обладает потенциальным ключом в качестве детерминанта (достигается за счет высокой независимости всех атрибутов друг от дргуа).

## Таблица applicant

- **id** - bigint, PRIMARY KEY
- **first_name** - text, NOT NULL, CHECK (length(first_name)<50),  имя работника может использоваться отдельно от фамилии
- **last_name** - text, NOT NULL, CHECK (length(last_name)<50), фамилия работника
- **city_id** - text, FK NOT NULL, связь с названием города
- **birth_date** - date, NOT NULL, дата рождения работника
- **path_to_profile_avatar** - text, NOT NULL, default (static/default_profile.png) аватарка для профиля работника
- **contacts** - text, CHECK (length(contacts)<150), контакты в которых уже сам работник должен указать название соцсети и ник
- **education** - text, CHECK (length(education)<150), образование, сам работник должен указать образование в той форме в которой он хочет
- **email** - text, UNIQUE, NOT NULL, CHECK (length(email)<50), почта работника UNIQUE
- **password_hash** - text, NOT NULL, CHECK (length(password_hash)<250), пароль работника
- **created_at** - timestamptz, NOT NULL, дата создания акаунта
- **updated_at** - timestamptz, NOT NULL, дата последнего изменения аккаунта

Relation **applicant to cv**:\
{id} -> applicant_id\
Связь: один работник ко многим резюме

## Таблица employer

- **id** - bigint, PRIMARY KEY
- **first_name** - text, NOT NULL, CHECK (length(first_name)<50), имя работника может использоваться отдельно от фамилии
- **last_name** - text, NOT NULL, CHECK (length(last_name)<50), фамилия работника
- **city_id** - text, FK, NOT NULL, связь с названием города
- **position** - text,  NOT NULL, CHECK (length(position)<50), должность занимаемая работодателем
- **company_name_id** - int, FK, NOT NULL, id названия компании
- **company_description** - text, NOT NULL, CHECK (length(company_description)<150), описании комании (той части за которую ответсвенен этот работодатель)
- **website** - text, NOT NULL, CHECK (length(website)<50), ссылка на сайт филиала компании за которую ответсвенен этот работодатель
- **path_to_profile_avatar** - text, NOT NULL, default (static/default_profile.png) аватарка для профиля работодателя
- **contacts** - text, CHECK (length(contacts)<50), контакты в которых уже сам работодатель должен указать название соцсети и ник
- **email** - text, UNIQUE, NOT NULL, CHECK (length(email)<50), почта работника UNIQUE
- **password_hash** - text, NOT NULL, CHECK (length(password_hash)<250), пароль работодателя
- **created_at** - timestamptz, NOT NULL, дата создания акаунта
- **updated_at** - timestamptz, NOT NULL, дата последнего изменения аккаунта

Relation **employer to vacancy**:\
{id} -> employer_id\
Связь: один работодатель ко многим вакансиям

## Таблица company

- **id** - int, PRIMARY KEY
- **name** - text, NOT NULL, UNIQUE, CHECK (length(name)<50), название комании

Relation **employer to company**:\
{company_name_id} -> id\
Связь: Много работодателей к одной компании

## Таблица city

- **id** - int, PRIMARY KEY
- **name** - text, UNIQUE, NOT NULL, CHECK (length(name)<50), название города
- **created_at** - timestamptz, NOT NULL, дата появления города в нашей бд
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления названия города

Relation **applicant to city**:\
{city_id} -> id\
Связь: Много работников к одному городу

Relation **employer to city**:\
{city_id} -> id\
Связь: Много работодателей к одному городу

## Таблица portfolio

- **id** - bigint, PRIMARY KEY
- **applicant_id** - bigint, FK, NOT NULL, id работника которому принадлежит это портфолио
- **name** - text, NOT NULL, CHECK (length(name)<50), название портфолио
- **created_at** - timestamptz, NOT NULL, дата появления города в нашей бд
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления названия города

Relation **applicant to portfolio**:\
{id} -> applicant_id\
Связь: Один работник ко многим портфолио

## Таблица applicant_creation_to_portfolio

- **portfolio_id** - bigint, FK, NOT NULL, id портфолио в котором будет это произведение
- **applicant_creation_id** - bigint, FK, NOT NULL, id произведения которое будет в этом портфолио
- **created_at** - timestamptz, NOT NULL, дата добавления произведения в портфолио
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления этого произведения в портфолио

Связь: многие произведения работодателей ко многим портфолио реализуемая с помощью промеждуточной таблицы applicant_creation_to_portfolio

Relation **portfolio to applicant_creation_to_portfolio**:\
{id} -> portfolio_id\
Связь: Одино портфолио может быть указано во многих строчках applicant_creation_to_portfolio\
(в одном портфолио много работ)

Relation **applicant_creation to applicant_creation_to_portfolio**:\
{id} -> applicant_creation_id\
Связь: Одино произведение работника может быть указано во многих строчках applicant_creation_to_portfolio\
(одна работа может отосится ко многим портфолио)

## Таблица cv_to_portfoli

- **cv_id** - bigint, FK, NOT NULL, id резюме к которому будет прикреплено это портфолио
- **portfolio_id** - bigint, FK, NOT NULL, id портфолио которое будет прикреплено к этому резюме
- **created_at** - timestamptz, NOT NULL, дата добавления портфолио к резюме
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления добавления резюме к портфолио

Связь: многие портфолио ко многим резюме реализуемая с помощью промеждуточной таблицы cv_to_portfoli

Relation **portfolio to cv_to_portfoli**:\
{id} -> portfolio_id\
Связь: Одино портфолио может быть указано во многих строчках cv_to_portfoli\
(одно портфолио можно прикрепить к нескольким резюме)

Relation **cv to cv_to_portfoli**:\
{id} -> cv_id\
Связь: Одино резюме может быть указано во многих строчках cv_to_portfoli\
(к одному резуюме можно прикрепить несколько портфолио)

## Таблица cv

- **id** - bigint, PRIMARY KEY
- **applicant_id** - bigint, FK, NOT NULL, id работника которому принадлежит это резюме
- **position_rus** - text, NOT NULL, CHECK (length(position_rus)<50), название желаемой должности на русском
- **position_eng** - text, CHECK (length(position_eng)<50), название желаемой должности на английском
- **job_search_status_id** - int, FK, NOT NULL, айди статуса поиска работы по этому резюме
- **working_experience** - text, CHECK (length(working_experience)<1000), описание опыта работы сотрудника
- **path_to_cv_avatar** - text, NOT NULL, default (static/default_profile.png) аватарка для резюме работника
- **created_at** - timestamptz, NOT NULL, дата создания резюме
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления резюме

## Таблица cv_subscriber

- **cv_id** - bigint, FK, NOT NULL, id резюме на которое подписался работодатель
- **employer_id** - bigint, FK, NOT NULL, id работодателя который подписался на это резюме
- **created_at** - timestamptz, NOT NULL, дата подписки работодателя на резюме
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления подписки работодателя на резюме

Relation **cv to cv_subscriber**:\
{id} -> cv_id\
Связь: Одно резюме может быть указано во многих строчках cv_subscriber\
(на одно резюме могут быть подписаны многие работодатели)

Relation **employer to cv_subscriber**:\
{id} -> cv_id\
Связь: Один работодатель может быть указан во многих строчках cv_subscriber\
(Один работодатель может быть подписан на несколько резюме)

## Таблица job_search_status

- **id** - int, PRIMARY KEY
- **status_name** - text, NOT NULL, UNIQUE, CHECK (length(status_name)<50), название статуса поиска работы
- **created_at** - timestamptz, NOT NULL, дата создания статуса поиска работы
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления статуса поиска работы

Relation **job_search_status to cv**:\
{id} -> job_search_status_id\
Связь: Один статус поиска работы может быть указзан во многих резюме\

## Таблица cv_to_creation_tags

- **creation_tag_id** - int, FK, NOT NULL, id  тега произведения который указан в этом резюме
- **cv_id** - bigint, FK, NOT NULL, id резюме в котором будет указан этот тег произведения
- **created_at** - timestamptz, NOT NULL, дата добавления тега произведения к резюме
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления тега произведения в резюме

Связь: многие теги произведений ко многим резюме реализуемая с помощью промеждуточной таблицы cv_to_creation_tags

Relation **cv to cv_to_creation_tags**:\
{id} -> cv_id\
Связь: Одино резюме может быть указано во многих строчках cv_to_creation_tags\
(Один тег произведения может быть указан во многих резюме)

Relation **creation_tags to cv_to_creation_tags**:\
{id} -> creation_tag_id\
Связь: Одиин тег произведения может быть указан во многих строчках cv_to_creation_tags\
(К одному резюме можно добавить несколько тегов произведений)

## Таблица applicant_creation_to_creation_tags

- **creation_tag_id** - int, FK, NOT NULL, id  тега для этого произведения
- **applicant_creation_id** - bigint, FK, NOT NULL, id произведения к которому будет прикреплен этот тег
- **created_at** - timestamptz, NOT NULL, дата добавления тега этому произведению
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления тега у этого произведения

Связь: многие теги произведений ко многим произведениям реализуемая с помощью промеждуточной таблицы applicant_creation_to_creation_tags

Relation **applicant_creation to applicant_creation_to_creation_tags**:\
{id} -> applicant_creation_id\
Связь: Одино произведение может быть указано во многих строчках applicant_creation_to_creation_tags\
(У одного произведения может быть много тегов)

Relation **creation_tags to applicant_creation_to_creation_tags**:\
{id} -> creation_tag_id\
Связь: Одиин тег произведения может быть указан во многих строчках applicant_creation_to_creation_tags\
(Один тег может быть отнесен ко многим произвденениям)

## Таблица vacancy_to_creation_tags

- **creation_tag_id** - int, FK, NOT NULL, id  тега произведения который указан в этой вакансии
- **vacancy_id** - bigint, FK, NOT NULL, id вакансии к которому будет прикреплено это портфолио
- **created_at** - timestamptz, NOT NULL, дата добавления тега произведения к вакансим
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления тега произведения в вакансии

Связь: многие теги произведений ко многим вакансиям реализуемая с помощью промеждуточной таблицы vacancy_to_creation_tags

Relation **vacancy to vacancy_to_creation_tags**:\
{id} -> vacancy_id\
Связь: Одина вакансия может быть указана во многих строчках vacancy_to_creation_tags\
(Одина вакансия может иметь много тегов желаемых произведений)

Relation **creation_tags to vacancy_to_creation_tags**:\
{id} -> creation_tag_id\
Связь: Одиин тег произведения может быть указан во многих строчках vacancy_to_creation_tags\
(Один тег произведения может быть указан у многих вакансий как желаемый)

## Таблица creation_tags

- **id** - int, PRIMARY KEY
- **creation_tag_name** - text, NOT NULL, CHECK (length(creation_tag_name)<50), UNIQUE, название тега произведения
- **created_at** - timestamptz, NOT NULL, дата добавления тега произведения
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления тега произведения

## Таблица vacancy

- **id** - bigint, PRIMARY KEY
- **employer_id** - bigint, FK, NOT NULL, id работодателя который разместил эту вакансию
- **salary** - int, NOT NULL зарботная плата предлагаемая сотруднику
- **position** - text, NOT NULL, CHECK (length(position)<50), должность предлагаемая сотруднику
- **description** - text, NOT NULL, CHECK (length(description)<100), описание вакансии от работодателя
- **work_type_id** - int, FK, NOT NULL, ссылка на тип работы (разовая, постоянная, пол ставки и тд)
- **path_to_company_avatar** - text, NOT NULL, default (static/default_company.png) логотип компании для вакансии
- **created_at** - timestamptz, NOT NULL, дата добавления вакансии
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления вакансии

## Таблица work_type

- **id** - int, PRIMARY KEY
- **work_type_name** - text, NOT NULL, CHECK (length(work_type_name)<50), тип работы (разовая, постоянная, пол ставки и тд)
- **created_at** - timestamptz, NOT NULL, дата подписки работника на вакансию
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления подписки работника на вакансию

Relation **work_type to vacancy**:\
{id} -> work_type_id\
Связь: Один тип работы может быть указан во многих вакансиях

## Таблица vacancy_subscriber

- **vacancy_id** - bigint, FK, NOT NULL, id вакансии на которую подписался работник
- **applicant_id** - bigint, FK, NOT NULL, id работника который подписался на эту вакансию
- **created_at** - timestamptz, NOT NULL, дата подписки работника на вакансию
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления подписки работника на вакансию

Связь: многие работники ко многим вакансиям реализуемая с помощью промеждуточной таблицы vacancy_subscriber

Relation **applicant to vacancy_subscriber**:\
{id} -> applicant_id\
Связь: Один работник может быть указан во многих строчках vacancy_subscriber\
(Один работник может быть подписан на несколько вакансий)

Relation **vacancy to vacancy_subscriber**:\
{id} -> vacancy_id\
Связь: Одна вакансия может быть указана во многих строчках vacancy_subscriber\
(На одну вакансию может быть подписано много работников)

## Таблица applicant_creation

- **id** - bigint, PRIMARY KEY
- **name** - text, NOT NULL, UNIQUE, CHECK (length(name)<50), название произвдения которое даст рабоник при загрузке произведения на наш сайт
- **creation** - text, NOT NULL, UNIQUE, адрес у нас а сервере где лежит эта работа
- **creation_type_id** - bigint, FK, NOT NULL, id типа произведения
- **created_at** - timestamptz, NOT NULL, дата выгрузки произведения на наш сайт
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления произведения на нашем сайте

## Таблица creation_type

- **id** - bigint, PRIMARY KEY
- **type_name** - text, NOT NULL, UNIQUE, CHECK (length(type_name)<50), название типа произведения
- **created_at** - timestamptz, NOT NULL, дата добавления типа произведения на наш сайт
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления типа произведения на наш сайт

Relation **creation_type to applicant_creation**:\
{id} -> creation_type_id\
Связь: Один тип произведения ко многим произведениям работников\

## Таблица applicant_session

- **id** - bigint, PRIMARY KEY
- **applicant_id** - bigint, FK, NOT NULL, id работника получившего эту сессию
- **session_token** - text, NOT NULL, UNIQUE, куки токен длинны 32
- **created_at** - timestamptz, NOT NULL, дата выдачи куки пользователю
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления куки пользователя

Relation **applicant to applicant_session**:\
{id} -> applicant_id\
Связь: У одного работника может быть несколько сессионных токенов для разных устройств

## Таблица employer_session

- **id** - bigint, PRIMARY KEY
- **employer_id** - bigint, FK, NOT NULL, id работодателя получившего эту сессию
- **session_token** - text, NOT NULL, UNIQUE, куки токен длинны 32
- **created_at** - timestamptz, NOT NULL, дата выдачи куки пользователю
- **updated_at** - timestamptz, NOT NULL, дата последнего обновления куки пользователя

Relation **employer to employer_session**:\
{id} -> employer_id\
Связь: У одного работодателя может быть несколько сессионных токенов для разных устройств

## Таблица applicant_rate_to_applicant_creation

- **rate** - int, рейтинг выставленный пользователем
- **applicant_id** - bigint, FK, NOT NULL, id работника поставившего рейтинг работе
- **applicant_creation_id** - bigint, FK, NOT NULL, id произведения которой работник поставил рейтинг
- **created_at** - timestamptz, NOT NULL, дата когда пользователь поставил рейтинг
- **updated_at** - timestamptz, NOT NULL, дата когда пользователь последний раз обновил рейтинг

Связь: многие работники ко многим произведениям работников (их рейтингу) реализуемая с помощью промеждуточной таблицы applicant_rate_to_applicant_creation

Relation **applicant to applicant_rate_to_applicant_creation**:\
{id} -> applicant_id\
Связь: Один работник может быть указан во многих строчках applicant_rate_to_applicant_creation\
(Один работник может поставить рейтинг разным произведениям)

Relation **applicant_creation to applicant_rate_to_applicant_creation**:\
{id} -> applicant_creation_id\
Связь: Одно произведение может быть указано во многих строчках applicant_rate_to_applicant_creation\
(У одной=го произведения может быть рейтинг от многих работников)

## Таблица employr_rate_to_applicant_creation

- **rate** - int, рейтинг выставленный пользователем
- **employr_id** - bigint, FK, NOT NULL, id работодателя поставившего рейтинг работе
- **applicant_creation_id** - bigint, FK, NOT NULL, id произведения которой работодатель поставил рейтинг
- **created_at** - timestamptz, NOT NULL, дата когда пользователь поставил рейтинг
- **updated_at** - timestamptz, NOT NULL, дата когда пользователь последний раз обновил рейтинг

Связь: многие работодатели ко многим произведениям работников (их рейтингу) реализуемая с помощью промеждуточной таблицы employr_rate_to_applicant_creation

Relation **employr to employr_rate_to_applicant_creation**:\
{id} -> employr_id\
Связь: Один работодатель может быть указан во многих строчках employr_rate_to_applicant_creation\
(Один работодатель может поставить рейтинг разным произведениям)

Relation **applicant_creation to employr_rate_to_applicant_creation**:\
{id} -> applicant_creation_id\
Связь: Одно произведение может быть указано во многих строчках employr_rate_to_applicant_creation\
(У одного произвдения может быть рейтинг от многих работодателей)
