UPDATE {{.Songs.Name}} SET {{.Songs.Fields.Fingerprinted}} = 1 WHERE {{.Songs.Fields.ID}}  = %s;