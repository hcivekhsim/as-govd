ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin'));

CREATE TABLE IF NOT EXISTS dictionaries (
    id SERIAL PRIMARY KEY,
    dictionary_type TEXT NOT NULL,
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (dictionary_type, code)
);

CREATE TABLE IF NOT EXISTS appeal_statuses (
    code TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    group_code TEXT NOT NULL,
    group_title TEXT NOT NULL,
    badge_type TEXT NOT NULL,
    is_final BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS appeals (
    id SERIAL PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    owner_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    applicant_category_code TEXT NOT NULL,
    topic_code TEXT NOT NULL,
    organization_name TEXT,
    appeal_kind_code TEXT NOT NULL,
    appeal_type_code TEXT NOT NULL,
    last_name TEXT NOT NULL,
    first_name TEXT NOT NULL,
    patronymic TEXT,
    identification_number TEXT NOT NULL,
    physical_status_code TEXT,
    social_category_code TEXT,
    email TEXT NOT NULL,
    is_abroad BOOLEAN NOT NULL DEFAULT false,
    citizenship_code TEXT,
    region TEXT,
    district TEXT,
    locality TEXT NOT NULL,
    street TEXT NOT NULL,
    house TEXT NOT NULL,
    apartment TEXT,
    note_address TEXT,
    postal_code TEXT,
    phone TEXT NOT NULL,
    content TEXT NOT NULL,
    planned_registration_date TIMESTAMP,
    current_status_code TEXT NOT NULL REFERENCES appeal_statuses(code),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT appeals_org_name_legal_only CHECK (
        (applicant_category_code = 'legal' AND organization_name IS NOT NULL)
        OR (applicant_category_code <> 'legal' AND organization_name IS NULL)
    )
);

CREATE TABLE IF NOT EXISTS appeal_files (
    id SERIAL PRIMARY KEY,
    appeal_id INT NOT NULL REFERENCES appeals(id) ON DELETE CASCADE,
    original_name TEXT NOT NULL,
    stored_name TEXT NOT NULL,
    storage_path TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (stored_name)
);

CREATE TABLE IF NOT EXISTS appeal_status_history (
    id SERIAL PRIMARY KEY,
    appeal_id INT NOT NULL REFERENCES appeals(id) ON DELETE CASCADE,
    from_status_code TEXT REFERENCES appeal_statuses(code),
    to_status_code TEXT NOT NULL REFERENCES appeal_statuses(code),
    comment TEXT,
    changed_by_user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_appeals_owner_user_id ON appeals(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_appeals_created_at ON appeals(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_appeals_status_code ON appeals(current_status_code);
CREATE INDEX IF NOT EXISTS idx_appeals_applicant_category_code ON appeals(applicant_category_code);
CREATE INDEX IF NOT EXISTS idx_appeals_topic_code ON appeals(topic_code);
CREATE INDEX IF NOT EXISTS idx_appeals_appeal_kind_code ON appeals(appeal_kind_code);
CREATE INDEX IF NOT EXISTS idx_appeals_appeal_type_code ON appeals(appeal_type_code);
CREATE INDEX IF NOT EXISTS idx_appeals_identification_number ON appeals(identification_number);
CREATE INDEX IF NOT EXISTS idx_appeals_planned_registration_date ON appeals(planned_registration_date);
CREATE INDEX IF NOT EXISTS idx_appeals_region ON appeals(region);
CREATE INDEX IF NOT EXISTS idx_appeals_district ON appeals(district);
CREATE INDEX IF NOT EXISTS idx_appeals_locality ON appeals(locality);
CREATE INDEX IF NOT EXISTS idx_appeals_street ON appeals(street);
CREATE INDEX IF NOT EXISTS idx_appeals_house ON appeals(house);
CREATE INDEX IF NOT EXISTS idx_appeals_postal_code ON appeals(postal_code);
CREATE INDEX IF NOT EXISTS idx_appeals_number ON appeals(number);

INSERT INTO dictionaries (dictionary_type, code, title, sort_order) VALUES
    ('applicant_category', 'individual', 'Физическое лицо', 10),
    ('applicant_category', 'legal', 'Юридическое лицо', 20),
    ('applicant_category', 'foreign', 'Иностранный гражданин', 30),
    ('applicant_category', 'stateless', 'Лицо без гражданства', 40),
    ('topic', 'general', 'Общий вопрос', 10),
    ('topic', 'housing', 'Жилищный вопрос', 20),
    ('topic', 'social_support', 'Социальная поддержка', 30),
    ('appeal_kind', 'application', 'Заявление', 10),
    ('appeal_kind', 'complaint', 'Жалоба', 20),
    ('appeal_kind', 'proposal', 'Предложение', 30),
    ('appeal_type', 'electronic', 'Электронное', 10),
    ('appeal_type', 'written', 'Письменное', 20),
    ('physical_status', 'adult', 'Совершеннолетний', 10),
    ('physical_status', 'retired', 'Пенсионер', 20),
    ('social_category', 'none', 'Без льготной категории', 10),
    ('social_category', 'veteran', 'Ветеран', 20),
    ('social_category', 'disabled', 'Лицо с инвалидностью', 30),
    ('citizenship', 'rb_citizen', 'Граждане Республики Беларусь', 10),
    ('citizenship', 'foreign_citizen', 'Иностранные граждане', 20)
ON CONFLICT (dictionary_type, code) DO NOTHING;

INSERT INTO appeal_statuses (code, title, group_code, group_title, badge_type, is_final) VALUES
    ('new', 'Новое', 'initial', 'Начальные', 'info', false),
    ('received', 'Поступило', 'initial', 'Начальные', 'info', false),
    ('registered', 'Зарегистрировано', 'initial', 'Начальные', 'info', false),
    ('in_progress', 'В работе', 'working', 'Рабочие', 'progress', false),
    ('under_review', 'На рассмотрении', 'working', 'Рабочие', 'progress', false),
    ('under_check', 'На проверке', 'working', 'Рабочие', 'progress', false),
    ('extended', 'Срок продлен', 'working', 'Рабочие', 'progress', false),
    ('rework', 'На доработке', 'working', 'Рабочие', 'progress', false),
    ('forwarded', 'Перенаправлено', 'intermediate', 'Промежуточные', 'warning', false),
    ('approval', 'На согласовании', 'intermediate', 'Промежуточные', 'warning', false),
    ('signature', 'На подписи', 'intermediate', 'Промежуточные', 'warning', false),
    ('awaiting_response', 'Ожидает ответа', 'intermediate', 'Промежуточные', 'warning', false),
    ('reviewed', 'Рассмотрено', 'final_positive', 'Финальные', 'success', true),
    ('answered', 'Ответ направлен', 'final_positive', 'Финальные', 'success', true),
    ('satisfied', 'Удовлетворено', 'final_positive', 'Финальные', 'success', true),
    ('approved', 'Одобрено', 'final_positive', 'Финальные', 'success', true),
    ('completed', 'Исполнено', 'final_positive', 'Финальные', 'success', true),
    ('rejected', 'Отказано', 'final_negative', 'Финальные отрицательные', 'warning', true)
ON CONFLICT (code) DO NOTHING;
