package noGen

import (
	"errors"
	"net/http"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

// GeneratePassNo generates passenger numbers for a list of passengers based on the plane's PassNoBase.
func GeneratePassNo(plane model.Plane, pass *[]model.Passenger) (err error) {
    if pass == nil {
        return errors.New("passenger list is nil")
    }

    baseValue := plane.PassNoBase

    p := model.Passenger{}
    err = databasehandler.Db.Unscoped().Where("pass_no BETWEEN ? AND ?", baseValue, baseValue + 99).Order("pass_no DESC").First(&p).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            err = nil
        } else {
            return
        }
    } else {
        baseValue = p.PassNo
    }



    for i := range *pass {
        baseValue++

        if baseValue < plane.PassNoBase + 100 {
            (*pass)[i].PassNo = baseValue
        } else {
            return cerror.New(http.StatusInternalServerError, "PASSENGER_NO_GEN", "Could not generate PassNo")
        }
    }

    return nil
}
