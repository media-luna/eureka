SELECT {{.Songs.Fields.Name}}, HEX({{.Songs.Fields.FileSHA1}}) AS {{.Songs.Fields.FileSHA1}}, {{.Songs.Fields.TotalHashes}}
FROM {{.Songs.Name}}
WHERE {{.Songs.Fields.ID}} = %s;