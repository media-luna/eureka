CREATE TABLE IF NOT EXISTS {{.Songs.Name}} (
    {{.Songs.Fields.ID}} SERIAL PRIMARY KEY,
    {{.Songs.Fields.Name}} VARCHAR(250) NOT NULL,
    {{.Songs.Fields.Fingerprinted}} BOOLEAN DEFAULT FALSE,
    {{.Songs.Fields.FileSHA1}} BYTEA NOT NULL,
    {{.Songs.Fields.TotalHashes}} INT NOT NULL DEFAULT 0,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);