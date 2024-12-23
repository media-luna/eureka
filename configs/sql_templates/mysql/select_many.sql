SELECT HEX({{.Fingerprints.Fields.Hash}}), {{.Songs.Fields.ID}}, {{.Fingerprints.Fields.Offset}}
FROM {{.Fingerprints.Name}}
WHERE {{.Fingerprints.Fields.Hash}} IN (%s);