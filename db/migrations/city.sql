CREATE TABLE IF NOT EXISTS public."city"
(
    id serial NOT NULL,
    city_name text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT city_primary_key PRIMARY KEY (id),
    CONSTRAINT city_city_name_length_check CHECK (length(city_name) <= 50) NOT VALID,
    CONSTRAINT city_city_name_unique UNIQUE (city_name),
)

---- create above / drop below ----

DROP TABLE IF EXISTS public."city";