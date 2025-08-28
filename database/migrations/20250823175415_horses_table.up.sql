CREATE TABLE IF NOT EXISTS horses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    date_of_birth DATE NOT NULL,
    gender INTEGER NOT NULL
);
