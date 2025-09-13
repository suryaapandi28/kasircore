-- baru bikin tabel providers
CREATE TABLE accounts_providers (
    f_kd_account UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    f_nama_account VARCHAR(150) NOT NULL,       
    f_email_account VARCHAR(150) UNIQUE NOT NULL,
    f_password TEXT NOT NULL,
    f_role_accout VARCHAR(150) NOT NULL,
    f_phone_account VARCHAR(20) NOT NULL,
    f_verification_account BOOLEAN NOT NULL DEFAULT FALSE,
    f_jwt_token TEXT,
    f_jwt_token_expired TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
