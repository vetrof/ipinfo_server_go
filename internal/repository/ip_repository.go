package repository

import (
	"database/sql"
	"ip_info_server/internal/models"
)

type IPRepository struct {
	db *sql.DB
}

func NewIPRepository(db *sql.DB) *IPRepository {
	return &IPRepository{db: db}
}

func (r *IPRepository) SaveIPInfo(info models.IPInfo) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO ip_info (
			ip, hostname, city, region, country,
			loc, org, postal, timezone, readme, user_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		info.IP, info.Hostname, info.City, info.Region, info.Country,
		info.Loc, info.Org, info.Postal, info.Timezone, info.Readme,
		info.UserID,
	)

	return err
}

func (r *IPRepository) GetHistoryByUserID(userID int) ([]models.IPInfo, error) {
	rows, err := r.db.Query(`
		SELECT ip, hostname, city, region, country,
		       loc, org, postal, timezone, readme, user_id
		FROM ip_info
		WHERE user_id = ?
		ORDER BY id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.IPInfo
	for rows.Next() {
		var info models.IPInfo
		err := rows.Scan(
			&info.IP, &info.Hostname, &info.City, &info.Region, &info.Country,
			&info.Loc, &info.Org, &info.Postal, &info.Timezone, &info.Readme, &info.UserID,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
