CREATE TABLE IF NOT EXISTS study_subjects (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       text NOT NULL,
    note_slug  text,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS study_sessions (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject_id       uuid NOT NULL REFERENCES study_subjects(id) ON DELETE CASCADE,
    session_date     date NOT NULL,
    duration_minutes int NOT NULL DEFAULT 0,
    topic            text NOT NULL DEFAULT '',
    created_at       timestamptz NOT NULL DEFAULT now()
);

INSERT INTO feature_flags (key, enabled, description) VALUES
    ('studies', true, 'Modulo de Estudios (materias y sesiones de estudio)'),
    ('studies.overview', true, 'Panorama de estudios: tiempo total y racha por materia (PRO)')
ON CONFLICT (key) DO NOTHING;
