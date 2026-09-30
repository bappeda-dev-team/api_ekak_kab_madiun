CREATE TABLE IF NOT EXISTS tb_tujuan_pemda_view (
    id              INT          NOT NULL AUTO_INCREMENT,
    id_tujuan_pemda INT          NOT NULL,
    is_hide         TINYINT(1)   NOT NULL DEFAULT 0,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_tujuan_pemda_view (id_tujuan_pemda)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
