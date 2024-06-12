package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome/model"
)

type loginStore struct {
	dbConn        *sql.DB
	encryptionKey string
}

func New(dbConn *sql.DB, encryptionKey string) Login {
	return &loginStore{
		dbConn:        dbConn,
		encryptionKey: encryptionKey,
	}
}

func (l loginStore) Get(filter *model.Filter) ([]model.Response, error) {
	var userLoginList []model.Response

	offset := 1
	getQuery := fmt.Sprintf("SELECT user_id, device_type, masked_ip, masked_device_id, locale, app_version, create_date FROM user_logins ORDER BY create_date DESC LIMIT %v OFFSET %v;", filter.Limit, filter.Page*offset)

	if filter.GroupDuplicates {
		getQuery = fmt.Sprintf("WITH DuplicateRecords AS (SELECT *, ROW_NUMBER() OVER (PARTITION BY masked_ip, masked_device_id ORDER BY create_date) AS rn FROM user_logins) SELECT user_id, device_type, masked_ip, masked_device_id, locale, app_version, create_date FROM DuplicateRecords WHERE rn > 1 ORDER BY create_date DESC LIMIT %v OFFSET %v;", filter.Limit, filter.Page*offset)
	}

	rows, err := l.dbConn.Query(getQuery)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error fetching records: %v", err.Error()))
	}

	for rows.Next() {
		var userLogin model.Response

		err = rows.Scan(&userLogin.UserID, &userLogin.DeviceType, &userLogin.IP, &userLogin.DeviceID, &userLogin.Locale, &userLogin.AppVersion, &userLogin.CreatedDate)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error fetching records: %v", err.Error()))
		}

		userLoginList = append(userLoginList, userLogin)
	}

	if !filter.IsEncrypted {
		for _, userLogin := range userLoginList {
			decryptedDeviceID, err := model.Decrypt(*userLogin.DeviceID, l.encryptionKey)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Error fetching records: %v", err.Error()))
			}

			decryptedIP, err := model.Decrypt(*userLogin.IP, l.encryptionKey)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Error fetching records: %v", err.Error()))
			}

			*userLogin.DeviceID = *decryptedDeviceID
			*userLogin.IP = *decryptedIP
		}
	}

	return userLoginList, nil
}
