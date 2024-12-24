SELECT COUNT({{.Songs.Fields.ID}}) AS n
FROM {{.Songs.Name}}
WHERE {{.Songs.Fields.Fingerprinted}} = 1;