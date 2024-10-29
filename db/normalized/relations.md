# Описание Базы Данных

## Таблица applicant

- **id** - bigint, PRIMARY KEY
- **first-name** - text, NOT NULL, CHECK (LENGTH(first-name)<50),  имя работника может использоваться отдельно от фамилии
- **last-name** - text, NOT NULL, CHECK (LENGTH(first-name)<50), фамилия работника
- **city-id** - text, FK NOT NULL, связь с названием города
- **birth-date** - date, NOT NULL, дата рождения работника
- **path-to-profile-avatar** - text, NOT NULL, default (static/default-profile.png) аватарка для профиля работника
- **contacts** - text, CHECK (LENGTH(first-name)<150), контакты в которых уже сам работник должен указать название соцсети и ник
- **education** - text, CHECK (LENGTH(first-name)<150), образование, сам работник должен указать образование в той форме в которой он хочет
- **email** - text, UNIQUE, NOT NULL, CHECK (LENGTH(first-name)<50), почта работника UNIQUE
- **password** - text, NOT NULL, CHECK (LENGTH(first-name)<50), пароль работника
- **created-at** - timestamptz, NOT NULL, дата создания акаунта
- **updated-at** - timestamptz, NOT NULL, дата последнего изменения аккаунта

Relation **applicant to cv**:\
{id} -> applicant-id\
Связь: один работник ко многим резюме

## Таблица employer

- **id** - bigint, PRIMARY KEY
- **first-name** - text, NOT NULL, CHECK (LENGTH(first-name)<50), имя работника может использоваться отдельно от фамилии
- **last-name** - text, NOT NULL, CHECK (LENGTH(first-name)<50), фамилия работника
- **city-id** - text, FK, NOT NULL, связь с названием города
- **position** - text,  NOT NULL, CHECK (LENGTH(first-name)<50), должность занимаемая работодателем
- **company-name-id** - int, FK, NOT NULL, id названия компании
- **company-description** - text, NOT NULL, CHECK (LENGTH(first-name)<150), описании комании (той части за которую ответсвенен этот работодатель)
- **website** - text, NOT NULL, CHECK (LENGTH(first-name)<50), ссылка на сайт филиала компании за которую ответсвенен этот работодатель
- **path-to-profile-avatar** - text, NOT NULL, default (static/default-profile.png) аватарка для профиля работодателя
- **contacts** - text, CHECK (LENGTH(first-name)<50), контакты в которых уже сам работодатель должен указать название соцсети и ник
- **email** - text, UNIQUE, NOT NULL, CHECK (LENGTH(first-name)<50), почта работника UNIQUE
- **password** - text, NOT NULL, CHECK (LENGTH(first-name)<50), пароль работодателя
- **created-at** - timestamptz, NOT NULL, дата создания акаунта
- **updated-at** - timestamptz, NOT NULL, дата последнего изменения аккаунта

Relation **employer to vacancy**:\
{id} -> employer-id\
Связь: один работодатель ко многим вакансиям

## Таблица company

- **id** - int, PRIMARY KEY
- **name** - text, NOT NULL, UNIQUE, CHECK (LENGTH(first-name)<50), название комании

Relation **employer to company**:\
{company-name-id} -> id\
Связь: Много работодателей к одной компании

## Таблица city

- **id** - int, PRIMARY KEY
- **name** - text, UNIQUE, NOT NULL, CHECK (LENGTH(first-name)<50), название города
- **created-at** - timestamptz, NOT NULL, дата появления города в нашей бд
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления названия города

Relation **applicant to city**:\
{city-id} -> id\
Связь: Много работников к одному городу

Relation **employer to city**:\
{city-id} -> id\
Связь: Много работодателей к одному городу

## Таблица portfolio

- **id** - bigint, PRIMARY KEY
- **applicant-id** - bigint, FK, NOT NULL, id работника которому принадлежит это портфолио
- **name** - text, NOT NULL, CHECK (LENGTH(first-name)<50), название портфолио
- **created-at** - timestamptz, NOT NULL, дата появления города в нашей бд
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления названия города

Relation **applicant to portfolio**:\
{id} -> applicant-id\
Связь: Один работник ко многим портфолио

## Таблица applicant-creation-to-portfolio

- **portfolio-id** - bigint, FK, NOT NULL, id портфолио в котором будет это произведение
- **applicant-creation-id** - bigint, FK, NOT NULL, id произведения которое будет в этом портфолио
- **created-at** - timestamptz, NOT NULL, дата добавления произведения в портфолио
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления этого произведения в портфолио

Связь: многие произведения работодателей ко многим портфолио реализуемая с помощью промеждуточной таблицы applicant-creation-to-portfolio

Relation **portfolio to applicant-creation-to-portfolio**:\
{id} -> portfolio-id\
Связь: Одино портфолио может быть указано во многих строчках applicant-creation-to-portfolio\
(в одном портфолио много работ)

Relation **applicant-creation to applicant-creation-to-portfolio**:\
{id} -> applicant-creation-id\
Связь: Одино произведение работника может быть указано во многих строчках applicant-creation-to-portfolio\
(одна работа может отосится ко многим портфолио)

## Таблица cv-to-portfoli

- **cv-id** - bigint, FK, NOT NULL, id резюме к которому будет прикреплено это портфолио
- **portfolio-id** - bigint, FK, NOT NULL, id портфолио которое будет прикреплено к этому резюме
- **created-at** - timestamptz, NOT NULL, дата добавления портфолио к резюме
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления добавления резюме к портфолио

Связь: многие портфолио ко многим резюме реализуемая с помощью промеждуточной таблицы cv-to-portfoli

Relation **portfolio to cv-to-portfoli**:\
{id} -> portfolio-id\
Связь: Одино портфолио может быть указано во многих строчках cv-to-portfoli\
(одно портфолио можно прикрепить к нескольким резюме)

Relation **cv to cv-to-portfoli**:\
{id} -> cv-id\
Связь: Одино резюме может быть указано во многих строчках cv-to-portfoli\
(к одному резуюме можно прикрепить несколько портфолио)

## Таблица cv

- **id** - bigint, PRIMARY KEY
- **applicant-id** - bigint, FK, NOT NULL, id работника которому принадлежит это резюме
- **position-rus** - text, NOT NULL, CHECK (LENGTH(first-name)<50), название желаемой должности на русском
- **position-eng** - text, CHECK (LENGTH(first-name)<50), название желаемой должности на английском
- **job-search-status-id** - int, FK, NOT NULL, айди статуса поиска работы по этому резюме
- **working-experience** - text, CHECK (LENGTH(first-name)<1000), описание опыта работы сотрудника
- **path-to-cv-avatar** - text, NOT NULL, default (static/default-profile.png) аватарка для резюме работника
- **created-at** - timestamptz, NOT NULL, дата создания резюме
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления резюме

## Таблица cv-subscriber

- **cv-id** - bigint, FK, NOT NULL, id резюме на которое подписался работодатель
- **employer-id** - bigint, FK, NOT NULL, id работодателя который подписался на это резюме
- **created-at** - timestamptz, NOT NULL, дата подписки работодателя на резюме
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления подписки работодателя на резюме

Relation **cv to cv-subscriber**:\
{id} -> cv-id\
Связь: Одно резюме может быть указано во многих строчках cv-subscriber\
(на одно резюме могут быть подписаны многие работодатели)

Relation **employer to cv-subscriber**:\
{id} -> cv-id\
Связь: Один работодатель может быть указан во многих строчках cv-subscriber\
(Один работодатель может быть подписан на несколько резюме)

## Таблица job-search-status

- **id** - int, PRIMARY KEY
- **status-name** - text, NOT NULL, UNIQUE, CHECK (LENGTH(first-name)<50), название статуса поиска работы
- **created-at** - timestamptz, NOT NULL, дата создания статуса поиска работы
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления статуса поиска работы

Relation **job-search-status to cv**:\
{id} -> job-search-status-id\
Связь: Один статус поиска работы может быть указзан во многих резюме\

## Таблица cv-to-creation-tags

- **creation-tag-id** - int, FK, NOT NULL, id  тега произведения который указан в этом резюме
- **cv-id** - bigint, FK, NOT NULL, id резюме в котором будет указан этот тег произведения
- **created-at** - timestamptz, NOT NULL, дата добавления тега произведения к резюме
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления тега произведения в резюме

Связь: многие теги произведений ко многим резюме реализуемая с помощью промеждуточной таблицы cv-to-creation-tags

Relation **cv to cv-to-creation-tags**:\
{id} -> cv-id\
Связь: Одино резюме может быть указано во многих строчках cv-to-creation-tags\
(Один тег произведения может быть указан во многих резюме)

Relation **creation-tags to cv-to-creation-tags**:\
{id} -> creation-tag-id\
Связь: Одиин тег произведения может быть указан во многих строчках cv-to-creation-tags\
(К одному резюме можно добавить несколько тегов произведений)

## Таблица applicant-creation-to-creation-tags

- **creation-tag-id** - int, FK, NOT NULL, id  тега для этого произведения
- **applicant-creation-id** - bigint, FK, NOT NULL, id произведения к которому будет прикреплен этот тег
- **created-at** - timestamptz, NOT NULL, дата добавления тега этому произведению
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления тега у этого произведения

Связь: многие теги произведений ко многим произведениям реализуемая с помощью промеждуточной таблицы applicant-creation-to-creation-tags

Relation **applicant-creation to applicant-creation-to-creation-tags**:\
{id} -> applicant-creation-id\
Связь: Одино произведение может быть указано во многих строчках applicant-creation-to-creation-tags\
(У одного произведения может быть много тегов)

Relation **creation-tags to applicant-creation-to-creation-tags**:\
{id} -> creation-tag-id\
Связь: Одиин тег произведения может быть указан во многих строчках applicant-creation-to-creation-tags\
(Один тег может быть отнесен ко многим произвденениям)

## Таблица vacancy-to-creation-tags

- **creation-tag-id** - int, FK, NOT NULL, id  тега произведения который указан в этой вакансии
- **vacancy-id** - bigint, FK, NOT NULL, id вакансии к которому будет прикреплено это портфолио
- **created-at** - timestamptz, NOT NULL, дата добавления тега произведения к вакансим
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления тега произведения в вакансии

Связь: многие теги произведений ко многим вакансиям реализуемая с помощью промеждуточной таблицы vacancy-to-creation-tags

Relation **vacancy to vacancy-to-creation-tags**:\
{id} -> vacancy-id\
Связь: Одина вакансия может быть указана во многих строчках vacancy-to-creation-tags\
(Одина вакансия может иметь много тегов желаемых произведений)

Relation **creation-tags to vacancy-to-creation-tags**:\
{id} -> creation-tag-id\
Связь: Одиин тег произведения может быть указан во многих строчках vacancy-to-creation-tags\
(Один тег произведения может быть указан у многих вакансий как желаемый)

## Таблица creation-tags

- **id** - int, PRIMARY KEY
- **creation-tag-name** - text, NOT NULL, CHECK (LENGTH(first-name)<50), UNIQUE, название тега произведения
- **created-at** - timestamptz, NOT NULL, дата добавления тега произведения
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления тега произведения

## Таблица vacancy

- **id** - bigint, PRIMARY KEY
- **employer-id** - bigint, FK, NOT NULL, id работодателя который разместил эту вакансию
- **salary** - int, NOT NULL зарботная плата предлагаемая сотруднику
- **position** - text, NOT NULL, CHECK (LENGTH(first-name)<50), должность предлагаемая сотруднику
- **description** - text, NOT NULL, CHECK (LENGTH(first-name)<100), описание вакансии от работодателя
- **work-type-id** - int, FK, NOT NULL, ссылка на тип работы (разовая, постоянная, пол ставки и тд)
- **path-to-company-avatar** - text, NOT NULL, default (static/default-company.png) логотип компании для вакансии
- **created-at** - timestamptz, NOT NULL, дата добавления вакансии
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления вакансии

## Таблица work-type

- **id** - int, PRIMARY KEY
- **work-type-name** - text, NOT NULL, CHECK (LENGTH(first-name)<50), тип работы (разовая, постоянная, пол ставки и тд)
- **created-at** - timestamptz, NOT NULL, дата подписки работника на вакансию
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления подписки работника на вакансию

Relation **work-type to vacancy**:\
{id} -> work-type-id\
Связь: Один тип работы может быть указан во многих вакансиях

## Таблица vacancy-subscriber

- **vacancy-id** - bigint, FK, NOT NULL, id вакансии на которую подписался работник
- **applicant-id** - bigint, FK, NOT NULL, id работника который подписался на эту вакансию
- **created-at** - timestamptz, NOT NULL, дата подписки работника на вакансию
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления подписки работника на вакансию

Связь: многие работники ко многим вакансиям реализуемая с помощью промеждуточной таблицы vacancy-subscriber

Relation **applicant to vacancy-subscriber**:\
{id} -> applicant-id\
Связь: Один работник может быть указан во многих строчках vacancy-subscriber\
(Один работник может быть подписан на несколько вакансий)

Relation **vacancy to vacancy-subscriber**:\
{id} -> vacancy-id\
Связь: Одна вакансия может быть указана во многих строчках vacancy-subscriber\
(На одну вакансию может быть подписано много работников)

## Таблица applicant-creation

- **id** - bigint, PRIMARY KEY
- **name** - text, NOT NULL, UNIQUE, CHECK (LENGTH(first-name)<50), название произвдения которое даст рабоник при загрузке произведения на наш сайт
- **creation** - text, NOT NULL, UNIQUE, адрес у нас а сервере где лежит эта работа
- **creation-type-id** - bigint, FK, NOT NULL, id типа произведения
- **created-at** - timestamptz, NOT NULL, дата выгрузки произведения на наш сайт
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления произведения на нашем сайте

## Таблица creation-type

- **id** - bigint, PRIMARY KEY
- **type-name** - text, NOT NULL, UNIQUE, CHECK (LENGTH(first-name)<50), название типа произведения
- **created-at** - timestamptz, NOT NULL, дата добавления типа произведения на наш сайт
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления типа произведения на наш сайт

Relation **creation-type to applicant-creation**:\
{id} -> creation-type-id\
Связь: Один тип произведения ко многим произведениям работников\

## Таблица applicant-session

- **id** - bigint, PRIMARY KEY
- **applicant-id** - bigint, FK, NOT NULL, id работника получившего эту сессию
- **cooky-token** - text, NOT NULL, UNIQUE, куки токен длинны 32
- **created-at** - timestamptz, NOT NULL, дата выдачи куки пользователю
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления куки пользователя

Relation **applicant to applicant-session**:\
{id} -> applicant-id\
Связь: У одного работника может быть несколько сессионных токенов для разных устройств

## Таблица employer-session

- **id** - bigint, PRIMARY KEY
- **employer-id** - bigint, FK, NOT NULL, id работодателя получившего эту сессию
- **cooky-token** - text, NOT NULL, UNIQUE, куки токен длинны 32
- **created-at** - timestamptz, NOT NULL, дата выдачи куки пользователю
- **updated-at** - timestamptz, NOT NULL, дата последнего обновления куки пользователя

Relation **employer to employer-session**:\
{id} -> employer-id\
Связь: У одного работодателя может быть несколько сессионных токенов для разных устройств

## Таблица applicant-rate-to-applicant-creation

- **rate** - int, рейтинг выставленный пользователем
- **applicant-id** - bigint, FK, NOT NULL, id работника поставившего рейтинг работе
- **applicant-creation-id** - bigint, FK, NOT NULL, id произведения которой работник поставил рейтинг
- **created-at** - timestamptz, NOT NULL, дата когда пользователь поставил рейтинг
- **updated-at** - timestamptz, NOT NULL, дата когда пользователь последний раз обновил рейтинг

Связь: многие работники ко многим произведениям работников (их рейтингу) реализуемая с помощью промеждуточной таблицы applicant-rate-to-applicant-creation

Relation **applicant to applicant-rate-to-applicant-creation**:\
{id} -> applicant-id\
Связь: Один работник может быть указан во многих строчках applicant-rate-to-applicant-creation\
(Один работник может поставить рейтинг разным произведениям)

Relation **applicant-creation to applicant-rate-to-applicant-creation**:\
{id} -> applicant-creation-id\
Связь: Одно произведение может быть указано во многих строчках applicant-rate-to-applicant-creation\
(У одной=го произведения может быть рейтинг от многих работников)

## Таблица employr-rate-to-applicant-creation

- **rate** - int, рейтинг выставленный пользователем
- **employr-id** - bigint, FK, NOT NULL, id работодателя поставившего рейтинг работе
- **applicant-creation-id** - bigint, FK, NOT NULL, id произведения которой работодатель поставил рейтинг
- **created-at** - timestamptz, NOT NULL, дата когда пользователь поставил рейтинг
- **updated-at** - timestamptz, NOT NULL, дата когда пользователь последний раз обновил рейтинг

Связь: многие работодатели ко многим произведениям работников (их рейтингу) реализуемая с помощью промеждуточной таблицы employr-rate-to-applicant-creation

Relation **employr to employr-rate-to-applicant-creation**:\
{id} -> employr-id\
Связь: Один работодатель может быть указан во многих строчках employr-rate-to-applicant-creation\
(Один работодатель может поставить рейтинг разным произведениям)

Relation **applicant-creation to employr-rate-to-applicant-creation**:\
{id} -> applicant-creation-id\
Связь: Одно произведение может быть указано во многих строчках employr-rate-to-applicant-creation\
(У одного произвдения может быть рейтинг от многих работодателей)
