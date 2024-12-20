# 2024_2_VKatuny
![Coverage](https://img.shields.io/badge/Coverage-60.3%25-yellow)

Данный репозиторий предназначен для хранения backend части проекта HeadHunter,
разрабатываемого командой VKатуны.

- [2024\_2\_VKatuny](#2024_2_vkatuny)
  - [Участники проекта](#участники-проекта)
    - [Члены команды](#члены-команды)
    - [Менторы](#менторы)
  - [О продукте](#о-продукте)

## Участники проекта

### Члены команды

- [Илья Андриянов](https://github.com/Regikon)
- [~~Виктория Гурьева~~](https://github.com/VikaGuryeva) (ушла из проекта)
- [Олег Музалев](https://github.com/Olgmuzalev13)
- [Михаил Черепнин](https://github.com/Ifelsik)

### Менторы

- **UI/UX**: Екатерина Гражданкина
- **Frontend**: Алексей Зотов
- **Backend**: Никита Архаров

## Ссылки на внешние ресурсы

- [Стандартная ссылка на деплой](https://uart.site)
- [Репозиторий фронтенда](https://github.com/frontend-park-mail-ru/2024_2_VKatuny)
- [Ссылка на документацию api (без рендера) в репозитории бэкенда](https://github.com/go-park-mail-ru/2024_2_VKatuny/tree/feature_vacancies-list-api/api)

## О продукте

Данный раздел будет заполнен по готовности четких представлений о продукте

## Полезные команды

### Линтер
Установка линтера
```
go install github.com/mgechev/revive@latest
```

```bash
make lint
```

### Тест

```bash
make tests
```

### Генерация swagger-файла
```bash
make api
```

### Генерация моков для тестов
```bash
make mock-gen
```

### Генерация easy-json
```bash
make easyjson-gen
```

## Необходимые библиотеки
[go-pdfium](https://github.com/klippa-app/go-pdfium?tab=readme-ov-file)  
[pdfium](https://github.com/bblanchon/pdfium-binaries/releases/download/chromium%2F6886/pdfium-linux-x64.tgz)  
  
[govips](https://github.com/davidbyttow/govips)
+
https://github.com/davidbyttow/govips/issues/100
(
sudo apt install pkg-config
sudo apt install libvips-dev
)
