INSERT IGNORE INTO {{.Fingerprints.Name}} (
        {{.Songs.Fields.ID}}
    ,   {{.Fingerprints.Fields.Hash}}
    ,   {{.Fingerprints.Fields.Offset}})
VALUES (%s, UNHEX(%s), %s);