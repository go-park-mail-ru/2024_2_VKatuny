# Описание Базы Данных

## Таблица applicant

- **id** - bigint, PRIMARY KEY
- **first-name** - text, NOT NULL имя работника может использоваться отдельно от фамилии
- **last-name** - text, NOT NULL фамилия работника
- **city-id** - text, NOT NULL связь с названием города
- **birth-date** - date, NOT NULL дата рождения работника
- **path-to-profile-avatar** - text, default "static/defaultprofile.png" аватарка для профиля работника
- **contacts** - text, контакты в которых уже сам работник должен указать название соцсети и ник
- **email** - text, NOT NULL почта работника UNIQUE
- **password** - text, NOT NULL пароль работника
- **created-at** - timestamptz, NOT NULL дата создания акаунта
- **updated-at** - timestamptz, NOT NULL дата последнего изменения аккаунта

Relation **Applicant to City**:

{city-id} -> id
Связь один город ко многим работникам

## Таблица City

- **id** - int, PRIMARY KEY
- **name** - text, UNIQUE NOT NULL название города
- **created-at** - timestamptz, дата появления города в нашей бд
- **updated-at** - timestamptz, дата последнего обновления названия города
