SELECT
    {{.Songs.Name}}
,   {{.Songs.Fields.Name}}
,   HEX({{.Songs.Fields.FileSHA1}}) AS {{.Songs.Fields.FileSHA1}}
,   {{.Songs.Fields.TotalHashes}}
,   `date_created`
FROM {{.Songs.Name}}
WHERE {{.Songs.Fields.Fingerprinted}} = 1;