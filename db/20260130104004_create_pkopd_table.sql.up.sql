CREATE TABLE pk_opd (
   id VARCHAR(255) NOT NULL,
   kode_opd VARCHAR(255) NOT NULL,
   nama_opd VARCHAR(255) NOT NULL,
   level_pk int NOT NULL,
   nip_atasan VARCHAR(255) NOT NULL,
   nama_atasan VARCHAR(255) NOT NULL,
   id_rekin_atasan VARCHAR(255) NOT NULL,
   rekin_atasan VARCHAR(255) NOT NULL,
   nip_pemilik_pk VARCHAR(255) NOT NULL,
   nama_pemilik_pk VARCHAR(255) NOT NULL,
   id_rekin_pemilik_pk VARCHAR(255) NOT NULL UNIQUE,
   rekin_pemilik_pk VARCHAR(255) NOT NULL,
   tahun int NOT NULL,
   keterangan VARCHAR(255),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
