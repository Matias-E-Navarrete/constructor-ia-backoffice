INSERT INTO feature_flags (key, enabled, description) VALUES
    ('habits.panorama', true, 'Panorama de habitos: heatmap, rachas y radar (PRO)')
ON CONFLICT (key) DO NOTHING;
