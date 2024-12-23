CREATE TABLE IF NOT EXISTS {{.Songs.Name}} (
    {{.Songs.Fields.ID}} MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
    {{.Songs.Fields.Name}} VARCHAR(250) NOT NULL,
    {{.Songs.Fields.Fingerprinted}} TINYINT DEFAULT 0,
    {{.Songs.Fields.FileSHA1}} BINARY(20) NOT NULL,
    {{.Songs.Fields.TotalHashes}} INT NOT NULL DEFAULT 0,
    date_created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_{{.Songs.Name}}_{{.Songs.Fields.ID}} PRIMARY KEY ({{.Songs.Fields.ID}}),
    CONSTRAINT uq_{{.Songs.Name}}_{{.Songs.Fields.ID}} UNIQUE KEY ({{.Songs.Fields.ID}})
) ENGINE=INNODB;