package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func db_open() (*sql.DB, error) {

	// get connection string from environment variables

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return db, err
	}

	return db, nil
}

func db_get_setting(setting_name string) string {

	setting_value := ""

	db, err := db_open()
	if err != nil {
		return ""
	}
	defer db.Close()

	sqlStatement := `select value from settings where name=$1;`

	row := db.QueryRow(sqlStatement, setting_name)
	err = row.Scan(&setting_value)
	if err != nil {
		return ""
	}

	return setting_value
}

func db_get_card_name_count(card_name string) (card_count int, err error) {

	card_count = 0

	db, err := db_open()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlStatement := `SELECT COUNT(card_id) FROM cards WHERE card_name = $1;`

	row := db.QueryRow(sqlStatement, card_name)
	err = row.Scan(&card_count)
	if err != nil {
		return 0, err
	}

	return card_count, nil
}

func db_insert_card(one_time_code string, k0_auth_key string, k2_cmac_key string, k3 string, k4 string,
	tx_max_sats int, day_max_sats int, lnurlw_enable bool, card_name string, uid_privacy bool,
	allow_neg_bal_ptr bool) error {

	lnurlw_enable_yn := "N"
	if lnurlw_enable {
		lnurlw_enable_yn = "Y"
	}

	uid_privacy_yn := "N"
	if uid_privacy {
		uid_privacy_yn = "Y"
	}

	allow_neg_bal_yn := "N"
	if allow_neg_bal_ptr {
		allow_neg_bal_yn = "Y"
	}

	db, err := db_open()
	if err != nil {
		return err
	}
	defer db.Close()

	// insert a new record into cards

	sqlStatement := `INSERT INTO cards` +
		` (one_time_code, k0_auth_key, k2_cmac_key, k3, k4, uid, last_counter_value,` +
		` lnurlw_request_timeout_sec, tx_limit_sats, day_limit_sats, lnurlw_enable,` +
		` one_time_code_used, card_name, uid_privacy, allow_negative_balance)` +
		` VALUES ($1, $2, $3, $4, $5, '', 0, 60, $6, $7, $8, 'N', $9, $10, $11);`
	res, err := db.Exec(sqlStatement, one_time_code, k0_auth_key, k2_cmac_key, k3, k4,
		tx_max_sats, day_max_sats, lnurlw_enable_yn, card_name, uid_privacy_yn,
		allow_neg_bal_yn)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("not one card record inserted")
	}

	return nil
}
