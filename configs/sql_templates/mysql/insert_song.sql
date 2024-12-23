INSERT INTO {{.Songs.Name}} ({{.Songs.Fields.Name}},{{.Songs.Fields.FileSHA1}},{{.Songs.Fields.TotalHashes}})
VALUES (%s, UNHEX(%s), %s);