CREATE TABLE IF NOT EXISTS workout_routines (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       text NOT NULL,
    exercises  text[] NOT NULL DEFAULT '{}',
    created_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO feature_flags (key, enabled, description) VALUES
    ('workouts.routines', true, 'Rutinas de entrenamiento reutilizables (PRO)')
ON CONFLICT (key) DO NOTHING;
