package main

import (
	"database/sql"
	"encoding/json"
	"net/url"

	_ "github.com/lib/pq"
)

type storage struct {
	backend *sql.DB
}

func openStorage(host, user, pass, db string) (storage, error) {
	v := url.Values{}
	v.Set("sslmode", "disable")
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, pass),
		Host:     host,
		Path:     "/" + url.PathEscape(db),
		RawQuery: v.Encode(),
	}
	dbconn, err := sql.Open("postgres", u.String())
	if err != nil {
		return storage{}, err
	}

	// force a connection now
	var tmp int
	if err := dbconn.QueryRow("SELECT 1;").Scan(&tmp); err != nil {
		dbconn.Close()
		return storage{}, err
	}

	return storage{dbconn}, nil
}

func (s storage) Close() error {
	return s.backend.Close()
}

type user struct {
	id   int
	name string
}

func (s storage) getUser(userStr, pass string) (user, error) {
	var uObj user
	if err := s.backend.QueryRow(
		`select "_user", "name" from "user" where "name"=$1 and "password"=$2`,
		userStr, pass,
	).Scan(&uObj.id, &uObj.name); err != nil {
		return uObj, err
	}
	return uObj, nil
}

func (s storage) getAgeByUserID(uid int) (int, error) {
	var age int
	if err := s.backend.QueryRow(
		`select "age" from "user" where "_user"=$1`,
		uid,
	).Scan(&age); err != nil {
		return 0, err
	}
	return age, nil
}

type coupon struct {
	cdataid int
	code    string
	data    map[string]string
}

func (s storage) getCouponForUser(uid int, couponName string) (coupon, error) {
	var cObj coupon
	var cid int
	if err := s.backend.QueryRow(
		`select "_coupon", "code" from "coupon" where "name"=$1`,
		couponName,
	).Scan(&cid, &cObj.code); err != nil {
		return cObj, err
	}

	// there is possible race condition here,
	// think what will happen, if this code called from 2 different go routine at same time
	// real code should use transaction or mutex
	var dataStr = "{}"
	if err := s.backend.QueryRow(
		`select "_coupon_data", "data" from "coupon_data" where "_user"=$1 and "_coupon"=$2`,
		uid, cid,
	).Scan(&cObj.cdataid, &dataStr); err != nil {
		if err != sql.ErrNoRows {
			return cObj, err
		}
		if err := s.backend.QueryRow(`
			insert into "coupon_data" ("_user", "_coupon", "data")
			values ($1, $2, $3)
			returning "_coupon_data"`,
			uid, cid, dataStr,
		).Scan(&cObj.cdataid); err != nil {
			return cObj, err
		}
	}
	if err := json.Unmarshal([]byte(dataStr), &cObj.data); err != nil {
		return cObj, err
	}
	return cObj, nil
}

func (s storage) saveCouponData(cdataid int, data map[string]string) error {
	dataBytes, _ := json.Marshal(data)
	dataStr := string(dataBytes)
	if _, err := s.backend.Exec(
		`update "coupon_data" set "data"=$2 where "_coupon_data"=$1`,
		cdataid, dataStr); err != nil {
		return err
	}
	return nil
}
