CREATE TABLE IF NOT EXISTS {{.Fingerprints.Name}} (
    {{.Fingerprints.Fields.Hash}} BINARY(10) NOT NULL,
    {{.Songs.Fields.ID}} MEDIUMINT UNSIGNED NOT NULL,
    {{.Fingerprints.Fields.Offset}} INT UNSIGNED NOT NULL,
    date_created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX ix_{{.Fingerprints.Name}}_{{.Fingerprints.Fields.Hash}} ({{.Fingerprints.Fields.Hash}}),
    CONSTRAINT uq_{{.Fingerprints.Name}}_{{.Songs.Fields.ID}}_{{.Fingerprints.Fields.Offset}}_{{.Fingerprints.Fields.Hash}}
        UNIQUE KEY ({{.Songs.Fields.ID}}, {{.Fingerprints.Fields.Offset}}, {{.Fingerprints.Fields.Hash}}),
    CONSTRAINT fk_{{.Fingerprints.Name}}_{{.Songs.Fields.ID}} FOREIGN KEY ({{.Songs.Fields.ID}})
        REFERENCES {{.Songs.Name}}({{.Songs.Fields.ID}}) ON DELETE CASCADE
) ENGINE=INNODB;