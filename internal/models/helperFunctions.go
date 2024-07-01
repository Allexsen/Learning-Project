package models

import "github.com/Allexsen/Learning-Project/internal/db"

func RetrieveAllRecordsByUserID(uid int64) ([]Record, error) {
	rows, err := db.DB.Query(`SELECT id, hours, minutes FROM practice_db.records WHERE user_id=?`, uid)
	if err != nil {
		return []Record{}, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		r := Record{UserID: uid}
		err = rows.Scan(&r.ID, &r.Hours, &r.Minutes)
		if err != nil {
			return []Record{}, err
		}

		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		return []Record{}, err
	}

	return records, nil
}
