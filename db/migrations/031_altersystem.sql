---- Write your migrate up statements here

ALTER SYSTEM SET log_line_prefix = '%t [%p]: [%l-1] '; -- префикс для журнала логов - время, тип (ошибка дебаг инфо) и текст
ALTER SYSTEM SET log_statement = 'all'; -- управляет тем, какие операторы SQL регистрируются
ALTER SYSTEM SET log_min_duration_statement = 1; -- 1мс для нашего небольшого количества данных уже долгий запрос
ALTER SYSTEM SET log_checkpoints = on; -- некоторые статистические данные включаются в сообщения журнала, включая количество записанных буферов и время, потраченное на их запись
ALTER SYSTEM SET log_connections = on; -- логгирует каждую попытку подключения к серверу, а также успешное завершения аутентификации и авторизации
ALTER SYSTEM SET log_disconnections = on; -- логгирует тоже что и log_connections только дял завершения соединения плюс продолжительность сеанса
ALTER SYSTEM SET log_lock_waits = on; -- управляет тем, будет ли создаваться сообщение журнала, когда сеанс ждет дольше deadlock_timeout для получения блокировки
ALTER SYSTEM SET log_temp_files = 0; -- управляет регистрацией имен и размеров временных файлов, я их не использую
ALTER SYSTEM SET log_autovacuum_min_duration = 0.1; -- заставляет регистрировать каждое действие, выполняемое автоочисткой, если оно выполнялось в течение как минимум указанного времени, автофункции должны выполнятся еще быстрее

ALTER SYSTEM SET logging_collector = true; -- этот параметр включает сборщик журналов , который является фоновым процессом, который захватывает сообщения журнала, отправленные в stderr , и перенаправляет их в файлы журналов
ALTER SYSTEM SET log_directory = '/data'; -- директория для logging_collector
ALTER SYSTEM SET log_filename = 'postgresql.log'; -- имя файла logging_collector
ALTER SYSTEM SET log_rotation_age = 1440; -- обновляю файл logging_collector каждые 24 часа
ALTER SYSTEM SET log_min_messages = 'INFO'; -- логирую все после уровня INFO
select pg_reload_conf(); -- необходимо перезагрузить конфиг чтобы применить изменения
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
