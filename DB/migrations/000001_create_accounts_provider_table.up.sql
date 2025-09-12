-- pertama bikin enum type untuk role
CREATE TYPE provider_role AS ENUM ('superadmin', 'admin', 'staff');

-- baru bikin tabel providers
CREATE TABLE accounts_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,       
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role provider_role DEFAULT 'admin',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
