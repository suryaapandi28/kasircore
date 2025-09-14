CREATE TABLE otp_verify (
    f_kd_otp UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    f_kd_account UUID NOT NULL,
    f_kode_otp VARCHAR(10) NOT NULL,
    f_otp_expired TIMESTAMPTZ NOT NULL,
    f_otp_via VARCHAR(50) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

ALTER TABLE otp_verify
ADD CONSTRAINT fk_account FOREIGN KEY (f_kd_account) REFERENCES accounts_providers(f_kd_account);
