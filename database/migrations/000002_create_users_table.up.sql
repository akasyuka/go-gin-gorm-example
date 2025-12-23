-- Создание таблицы users
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,          -- ID соответствует uint + primaryKey
                       email TEXT NOT NULL UNIQUE,        -- Email уникальный и not null
                       name TEXT NOT NULL,                -- Name not null
                       created_at TIMESTAMP NOT NULL DEFAULT now()  -- CreatedAt с текущей датой
);
