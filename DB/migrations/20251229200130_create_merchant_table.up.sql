CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE merchants (
    -- PRIMARY KEY
    f_kode_merchant UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- IDENTITAS
    f_nama_merchant VARCHAR(150) NOT NULL,
    f_jenis_usaha VARCHAR(50) NOT NULL, -- F&B, Retail, Jasa

    -- KONTAK
    f_email_merchant VARCHAR(150),
    f_phone_merchant VARCHAR(30),

    -- LOKASI
    f_alamat_merchant TEXT,
    f_kota VARCHAR(100),
    f_provinsi VARCHAR(100),
    f_kode_pos VARCHAR(10),

    -- PENGATURAN POS
    f_currency VARCHAR(10) DEFAULT 'IDR',
    f_ppn_enabled BOOLEAN DEFAULT FALSE,
    f_ppn_persen NUMERIC(5,2) DEFAULT 11.00,

    -- STATUS
    f_status_merchant BOOLEAN DEFAULT TRUE,

    -- AUDIT
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
