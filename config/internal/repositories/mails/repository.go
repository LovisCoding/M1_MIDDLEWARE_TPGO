package mails

import (
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
)

func GetAll() ([]models.Mail, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	rows, err := db.Query("SELECT id, mail FROM Mail")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mails []models.Mail
	for rows.Next() {
		var m models.Mail
		if err := rows.Scan(&m.ID, &m.Mail); err != nil {
			return nil, err
		}
		mails = append(mails, m)
	}
	return mails, nil
}

func GetById(id int64) (*models.Mail, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	var m models.Mail
	err = db.QueryRow("SELECT id, mail FROM Mail WHERE id = ?", id).Scan(&m.ID, &m.Mail)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func Create(mail string) (int64, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return 0, err
	}
	defer helpers.CloseDB(db)

	res, err := db.Exec("INSERT INTO Mail (mail) VALUES (?)", mail)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Delete(id int64) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("DELETE FROM Mail WHERE id = ?", id)
	return err
}

func Update(id int64, mail string) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("UPDATE Mail SET mail = ? WHERE id = ?", mail, id)
	return err
}
