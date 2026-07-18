CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email         text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    plan          text NOT NULL DEFAULT 'free' CHECK (plan IN ('free', 'pro')),
    role          text NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS groups (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name          text NOT NULL,
    kind          text NOT NULL CHECK (kind IN ('family', 'coaching')),
    owner_user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS group_members (
    group_id  uuid NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id   uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role      text NOT NULL CHECK (role IN ('owner', 'member', 'coach', 'client')),
    joined_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE IF NOT EXISTS notes (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug          text NOT NULL,
    title         text NOT NULL,
    body_markdown text NOT NULL DEFAULT '',
    kind          text NOT NULL DEFAULT 'freeform' CHECK (kind IN ('freeform', 'profile')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, slug)
);

CREATE TABLE IF NOT EXISTS habits (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        text NOT NULL,
    note_slug   text,
    created_at  timestamptz NOT NULL DEFAULT now(),
    archived_at timestamptz
);

CREATE TABLE IF NOT EXISTS habit_logs (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id   uuid NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    log_date   date NOT NULL,
    completed  boolean NOT NULL DEFAULT true,
    UNIQUE (habit_id, log_date)
);

CREATE TABLE IF NOT EXISTS workout_sessions (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_date date NOT NULL,
    notes        text NOT NULL DEFAULT '',
    note_slug    text,
    created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workout_sets (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id    uuid NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    exercise_name text NOT NULL,
    set_number    int NOT NULL,
    reps          int NOT NULL,
    weight_kg     numeric(7, 2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        text NOT NULL CHECK (type IN ('income', 'expense')),
    amount      numeric(12, 2) NOT NULL,
    category    text NOT NULL DEFAULT 'general',
    description text NOT NULL DEFAULT '',
    tx_date     date NOT NULL,
    note_slug   text,
    group_id    uuid REFERENCES groups(id) ON DELETE SET NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS feature_flags (
    key         text PRIMARY KEY,
    enabled     boolean NOT NULL DEFAULT true,
    description text NOT NULL DEFAULT '',
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS roadmap_items (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title       text NOT NULL,
    description text NOT NULL DEFAULT '',
    kind        text NOT NULL DEFAULT 'task' CHECK (kind IN ('feature', 'bug', 'task')),
    status      text NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'in_progress', 'done')),
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

INSERT INTO feature_flags (key, enabled, description) VALUES
    ('habits', true, 'Modulo de Habitos'),
    ('workouts', true, 'Modulo de Entrenamientos'),
    ('finance', true, 'Modulo de Finanzas'),
    ('notes', true, 'Modulo de Notas (vault estilo Obsidian)'),
    ('groups', true, 'Cuentas familiares y coaching'),
    ('habits.stats', true, 'Estadisticas de habitos (PRO)'),
    ('workouts.progress', true, 'Progreso de entrenamientos (PRO)'),
    ('finance.summary', true, 'Resumen financiero (PRO)'),
    ('finance.export', true, 'Exportar finanzas a CSV (PRO)'),
    ('notes.graph', true, 'Vista de grafo de notas (PRO)'),
    ('notes.export', true, 'Exportar vault de Obsidian (PRO)')
ON CONFLICT (key) DO NOTHING;
