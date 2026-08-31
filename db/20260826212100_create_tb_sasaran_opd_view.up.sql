CREATE TABLE IF NOT EXISTS tb_sasaran_opd_view (
    id         INT          NOT NULL AUTO_INCREMENT,
    id_pokin   INT          NOT NULL,
    is_hide    TINYINT(1)   NOT NULL DEFAULT 0,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_sasaran_opd_view_pokin (id_pokin)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
