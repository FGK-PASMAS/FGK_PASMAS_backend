package debugservice

import dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"

func TruncateData() {
    dh.Db.Exec("TRUNCATE TABLE passengers RESTART IDENTITY")
    dh.Db.Exec("TRUNCATE TABLE division RESTART IDENTITY")
}
