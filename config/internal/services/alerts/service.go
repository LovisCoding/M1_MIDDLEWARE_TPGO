package alerts

import (
	"middleware/config/internal/models"
	repository "middleware/config/internal/repositories/alerts"

	"github.com/sirupsen/logrus"
)

func GetAlertsForResource(resId int64) ([]models.Alert, error) {
	alerts, err := repository.GetAlertsForResource(resId)
	if err != nil {
		logrus.Errorf("error retrieving alerts for resource %d : %s", resId, err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving alerts",
		}
	}
	return alerts, nil
}

func AddAlert(mailId, resId int64) error {
	err := repository.AddAlert(mailId, resId)
	if err != nil {
		logrus.Errorf("error adding alert (mail:%d, res:%d) : %s", mailId, resId, err.Error())
		return &models.ErrorGeneric{
			Message: "Something went wrong while adding alert",
		}
	}
	return nil
}

func RemoveAlert(mailId, resId int64) error {
	err := repository.RemoveAlert(mailId, resId)
	if err != nil {
		logrus.Errorf("error removing alert (mail:%d, res:%d) : %s", mailId, resId, err.Error())
		return &models.ErrorGeneric{
			Message: "Something went wrong while removing alert",
		}
	}
	return nil
}
