package pasmasservice_test

import (
	"testing"

	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB(t *testing.T) {
    dsn := databasehandler.GetConnectionString()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{

    })
    if err != nil {
        assert.Nil(t, err)
        t.FailNow()
    } else  {
        databasehandler.InitGorm(db)
    }
}
func TestGetDivisions(t *testing.T) {
    initDB(t)

    div, err := pasmasservice.GetDivisions()
    if err != nil {
        assert.Nil(t, err)
        t.FailNow()
    }

    assert.IsType(t, []model.Division{}, div)
    assert.GreaterOrEqual(t, 3, len(div))
}
