package resources

import (
	"database/sql"
	"fmt"
	"middleware/config/internal/models"
	repository "middleware/config/internal/repositories/resources"

	"github.com/sirupsen/logrus"
)

func GetAllResources() ([]models.Ressource, error) {
	res, err := repository.GetAll()
	if err != nil {
		logrus.Errorf("error retrieving resources : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving resources",
		}
	}
	return res, nil
}

func GetResourceById(id int64) (*models.Ressource, error) {
	res, err := repository.GetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.ErrorNotFound{
				Message: "resource not found",
			}
		}
		logrus.Errorf("error retrieving resource %d : %s", id, err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving resource %d", id),
		}
	}
	return res, nil
}

func CreateResource(id int64, name string, resType *string) (int64, error) {
	newId, err := repository.Create(id, name, resType)
	if err != nil {
		logrus.Errorf("error creating resource %s : %s", name, err.Error())
		return 0, &models.ErrorGeneric{
			Message: "Something went wrong while creating resource",
		}
	}
	return newId, nil
}

func UpdateResource(id int64, name string, resType *string) error {
	err := repository.Update(id, name, resType)
	if err != nil {
		logrus.Errorf("error updating resource %d : %s", id, err.Error())
		return &models.ErrorGeneric{
			Message: "Something went wrong while updating resource",
		}
	}
	return nil
}

func DeleteResource(id int64) error {
	err := repository.Delete(id)
	if err != nil {
		logrus.Errorf("error deleting resource %d : %s", id, err.Error())
		return &models.ErrorGeneric{
			Message: "Something went wrong while deleting resource",
		}
	}
	return nil
}
