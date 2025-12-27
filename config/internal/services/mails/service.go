package mails

import (
	"database/sql"
	"fmt"
	"middleware/config/internal/models"
	repository "middleware/config/internal/repositories/mails"

	"github.com/sirupsen/logrus"
)

func GetAllMails() ([]models.Mail, error) {
	mails, err := repository.GetAll()
	if err != nil {
		logrus.Errorf("error retrieving mails : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving mails",
		}
	}
	return mails, nil
}

func GetMailById(id int64) (*models.Mail, error) {
	mail, err := repository.GetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.ErrorNotFound{
				Message: "mail not found",
			}
		}
		logrus.Errorf("error retrieving mail %d : %s", id, err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving mail %d", id),
		}
	}
	return mail, nil
}

func CreateMail(mailStr string) (int64, error) {
	id, err := repository.Create(mailStr)
	if err != nil {
		logrus.Errorf("error creating mail %s : %s", mailStr, err.Error())
		return 0, &models.ErrorGeneric{
			Message: "Something went wrong while creating mail",
		}
	}
	return id, nil
}

func UpdateMail(id int64, mailStr string) error {
	err := repository.Update(id, mailStr)
	if err != nil {
		logrus.Errorf("error updating mail %d : %s", id, err.Error())
		return &models.ErrorGeneric{
			Message: "Something went wrong while updating mail",
		}
	}
	return nil
}

func DeleteMail(id int64) error {
	err := repository.Delete(id)
	if err != nil {
		logrus.Errorf("error deleting mail %d : %s", id, err.Error())
		return &models.ErrorGeneric{
			Message: "Something went wrong while deleting mail",
		}
	}
	return nil
}
