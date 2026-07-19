CREATE TABLE IF NOT EXISTS assistant_messages (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       text NOT NULL CHECK (role IN ('user', 'assistant')),
    content    text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO feature_flags (key, enabled, description) VALUES
    ('assistant', true, 'Asistente de IA (PRO, requiere ANTHROPIC_API_KEY)')
ON CONFLICT (key) DO NOTHING;
