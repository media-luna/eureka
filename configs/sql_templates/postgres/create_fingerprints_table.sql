CREATE TABLE IF NOT EXISTS {{.Fingerprints.Name}} (
    {{.Fingerprints.Fields.Hash}} BYTEA NOT NULL,
    {{.Songs.Fields.ID}} INT NOT NULL,
    {{.Songs.Fields.Name}} INT NOT NULL,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    UNIQUE ({{.Songs.Fields.ID}}, {{.Songs.Fields.Name}}, {{.Fingerprints.Fields.Hash}}),
    FOREIGN KEY ({{.Songs.Fields.ID}}) REFERENCES {{.Songs.Name}}({{.Songs.Fields.ID}}) ON DELETE CASCADE
);