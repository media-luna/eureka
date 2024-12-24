SELECT {{.Songs.Fields.ID}}, {{.Fingerprints.Fields.Offset}} 
FROM {{.Fingerprints.Name}}
WHERE {{.Fingerprints.Fields.Hash}} = UNHEX(%s);