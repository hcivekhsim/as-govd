-- ЗАПУСТИ В БД
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, 
    user_type TEXT NOT NULL CHECK(user_type IN ('individual', 'legal', 'foreign', 'stateless')),
    name_organization TEXT, 
    login VARCHAR(20) NOT NULL UNIQUE, 
    password CHAR(60) NOT NULL, 
    email TEXT UNIQUE,
    phone VARCHAR(20) UNIQUE,
    last_name TEXT,
    first_name TEXT,
    patronymic TEXT,
    is_abroad BOOLEAN DEFAULT false,
    locality TEXT DEFAULT 'г. Барановичи',
    street TEXT NOT NULL,
    house TEXT NOT NULL,
    apartment TEXT,
    note_address TEXT,
    postal_code VARCHAR(10),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT check_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    CONSTRAINT check_phone CHECK (phone ~ '^\+[0-9]{10,15}$'),
    -- CONSTRAINT check_last_name CHECK (last_name ~ '^([А-ЯЁ][а-яё]+|[A-Z][a-z]+)$'),
    -- CONSTRAINT check_first_name CHECK (first_name ~ '^([А-ЯЁ][а-яё]+|[A-Z][a-z]+)$'),
    -- CONSTRAINT check_patronymic CHECK (patronymic ~ '^([А-ЯЁ][а-яё]+|[A-Z][a-z]+)$'),
    CONSTRAINT check_name_organization CHECK ((user_type = 'legal' AND name_organization IS NOT NULL) OR (user_type <> 'legal' AND name_organization IS NULL))
);

CREATE TABLE IF NOT EXISTS users_auth(
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  session_token_hash TEXT NOT NULL UNIQUE,
  csrf_token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL, 
  created_at TIMESTAMP NOT NULL DEFAULT now()
);




-- 1 проверка пароль: от 8 символов обязательно минимум 1 цифру и латинские буквы в верхнем и нижнем регистре BACKEND
-- 2 проверка имя, фамилия, отчество должно начинаться с большой буквы дальше маленькими [+]
-- 3 проверка номера телефона: + обязателен, от 10 до 15 цифр [+]
-- 4 проверка валидность email [+]
