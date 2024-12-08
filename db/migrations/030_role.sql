-- Write your migrate up statements here
CREATE ROLE vkatun WITH LOGIN PASSWORD 'passIMOapp' CONNECTION LIMIT 1; -- подключение только одно
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO vkatun;-- Приложение не создает таблицы поэтому на это у него прав нет
ALTER ROLE vkatun SET statement_timeout = '30s'; -- Ожидание, которое будет еще не очень критично пользователю, но может обеспечить нам защиту от DOS атак
ALTER ROLE vkatun SET lock_timeout = '10s'; -- Параметр имеет смысл только если он меньше statement_timeout иначе statement_timeout сработает раньше

LOAD 'auto_explain';
SET auto_explain.log_min_duration = 1; -- время выполнения оператора, в миллисекундах, при превышении которого план оператора будет протоколироваться у нас мало данных поэтому 1 миллисекунда это уже долгий запрос
SET auto_explain.log_analyze = true; -- в протокол будет записываться вывод команды EXPLAIN ANALYZE, а не простой EXPLAIN (более полная информация)
SET auto_explain.log_triggers = true; -- в протокол будет записываться статистика выполнения триггеров
SET auto_explain.log_nested_statements = true; -- протоколированию могут подлежать и вложенные операторы (операторы, выполняемые внутри функции)

--CREATE extension pg_profile schema public;
CREATE EXTENSION pg_stat_statements;
SET pg_stat_statements.max = 10000; -- максимальное число операторов, отслеживаемых модулем (не хотелось бы терять статистику при высокой нагрузке поэтому поставил в 2 раза выше дефолтного значения)

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
