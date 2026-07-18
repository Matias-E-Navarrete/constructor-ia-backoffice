CREATE TABLE IF NOT EXISTS tasks (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title         text NOT NULL,
    description   text NOT NULL DEFAULT '',
    category      text NOT NULL DEFAULT 'General',
    subcategory   text,
    due_date      date,
    due_time      text,
    repeat_rule   text NOT NULL DEFAULT 'none' CHECK (repeat_rule IN ('none', 'daily', 'weekly', 'monthly')),
    priority      text NOT NULL DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    status        text NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'in_progress', 'done')),
    quadrant      text CHECK (quadrant IN ('do_now', 'schedule', 'delegate', 'eliminate')),
    sort_order    int NOT NULL DEFAULT 0,
    scheduled_at  timestamptz,
    subtasks      jsonb NOT NULL DEFAULT '[]',
    note_slug     text,
    completed_at  timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS inbox_items (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    text NOT NULL,
    pinned     boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            text NOT NULL,
    type            text NOT NULL DEFAULT 'Cuenta',
    currency        text NOT NULL DEFAULT 'USD',
    initial_balance numeric(14, 2) NOT NULL DEFAULT 0,
    active          boolean NOT NULL DEFAULT true,
    created_at      timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS account_id uuid REFERENCES accounts(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS currency text NOT NULL DEFAULT 'USD';
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS exchange_rate numeric(14, 6);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS installments_total int;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS installment_number int;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS recurring boolean NOT NULL DEFAULT false;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS recurrence_interval text;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS method text NOT NULL DEFAULT '';

ALTER TABLE users ADD COLUMN IF NOT EXISTS planner_start_hour int NOT NULL DEFAULT 6;
ALTER TABLE users ADD COLUMN IF NOT EXISTS planner_end_hour int NOT NULL DEFAULT 23;

INSERT INTO feature_flags (key, enabled, description) VALUES
    ('tasks', true, 'Modulo de Tareas (Eisenhower, Kanban, Planner)'),
    ('inbox', true, 'Bandeja de entrada GTD')
ON CONFLICT (key) DO NOTHING;
