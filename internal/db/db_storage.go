package databasestorage

import (
	"database/sql"
	"fmt"

	bn "github.com/AlexGrtsnk/go-prof-diploma-2/internal/bindata"

	flw "github.com/AlexGrtsnk/go-prof-diploma-2/internal/json_parser"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const drriver = "sqlite3"
const dbbName = "shortenerdbs.db"

// NewDB создает новую sql сущность с конкретным необходимым нам драйвером
func NewDB() (*sql.DB, error) {
	dbname, driverTemp, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	sqliteDB, err := sql.Open(driverTemp, dbname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sqlite DB")
	}

	return sqliteDB, nil
}

// RunMigrateScripts запускает скрипты миграции в зависимости от типа драйвера
func RunMigrateScripts(db *sql.DB) error {
	var driver database.Driver
	var err error
	dbNameTemp, _, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	if dbNameTemp == dbbName {
		driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
	} else {
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	}
	if err != nil {
		return fmt.Errorf("creating db driver failed %s", err)
	}

	res := bindata.Resource(bn.AssetNames(),
		func(name string) ([]byte, error) {
			return bn.Asset(name)
		})

	d, _ := bindata.WithInstance(res)
	m, err := migrate.NewWithInstance("go-bindata", d, dbNameTemp, driver)
	if err != nil {
		return fmt.Errorf("initializing db migration failed %s", err)
	}
	if dbNameTemp == dbbName {
		_ = m.Steps(-1)
		err = m.Steps(1)
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating database failed %s", err)
		}
	} else {
		_ = m.Down()
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating database failed %s", err)
		}
	}
	return nil
}

// DataBaseCfg конфигурирует базу данных
func DataBaseCfg(flagRunAddr string, apiRunAddr string, fileName string, trustedSubnet string) (err error) {
	db, err := NewDB()
	if err != nil {
		return err
	}

	defer db.Close()
	err = RunMigrateScripts(db)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO cfg(flagRunAddr, apiRunAddr, flnm, tssb) VALUES ('` + string(flagRunAddr) + `', '` + string(apiRunAddr) + `', '` + fileName + `', '` + trustedSubnet + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

// DataBasePingHandler хендлер для проверки отклика базы данных
func DataBasePingHandler() (err error) {
	_, driverTemp, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	dbName := fmt.Sprintf("host=%s port=%s  user=%s password=%s dbname=%s sslmode=disable",
		`postgres`, `5432`, `postgres`, `postgres`, `praktikum`)
	err = DataBasePing(dbName, driverTemp)
	if err != nil {
		err = DataBaseSelfConfigUpdate(dbbName, drriver)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

// DataBasePing проверяет, отвечает ли наша база данных
func DataBasePing(dbbname string, driver string) (err error) {
	var db *sql.DB
	var res string
	db, err = sql.Open(driver, dbbname)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "SELECT 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}
	rows.Next()
	err = rows.Scan(&res)
	if err != nil {
		return err
	}
	return nil
}

// DataBaseAPIAddressSelect возвращает адрес апи
func DataBaseAPIAddressSelect() (apiAddress string, err error) {
	var db *sql.DB
	var apiRunAddr string
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return "", err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return "", err
	}
	defer db.Close()
	quer := "SELECT flnm FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", rows.Err()
	}
	rows.Next()
	err = rows.Scan(&apiRunAddr)
	if err != nil {
		return "", err
	}
	return apiRunAddr, nil
}

// DataBaseStartConfig конфигурирует изначаотный конфиг для записи
func DataBaseStartConfig(dbName string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return err
	}
	defer db.Close()
	var driver string
	if dbName != "localhost" {
		driver = "pgx"
	} else {
		driver = "sqlite3"
		dbName = dbbName
	}
	sts1 := `
	DROP TABLE IF EXISTS cfg;
	CREATE TABLE cfg (id INTEGER PRIMARY KEY, dbbname TEXT, driver TEXT);
	INSERT INTO cfg(dbbname, driver) VALUES ('` + string(dbName) + `', '` + string(driver) + `');`
	_, err = db.Exec(sts1)

	if err != nil {
		return err
	}
	return nil
}

// DataBaseSelfConfigGet возвращает конфиг, на котором сейчас работает база данных
func DataBaseSelfConfigGet() (dbbname string, driver string, err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return "", "", err
	}
	quer := "SELECT dbbname FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", "", rows.Err()
	}
	rows.Next()
	var dbNameTemp string
	err = rows.Scan(&dbNameTemp)
	if err != nil {
		return "", "", err
	}
	quer = "SELECT driver FROM cfg WHERE id = 1;"
	rows, err = db.Query(quer)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", "", rows.Err()
	}
	rows.Next()
	var driverTemp string
	err = rows.Scan(&driverTemp)
	if err != nil {
		return "", "", err
	}
	return dbNameTemp, driverTemp, nil
}

// DataBaseSelfConfigUpdate изменяет конфиг, записанный в базу данных
func DataBaseSelfConfigUpdate(dbbname string, driver string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "UPDATE cfg SET dbbname='" + dbbname + "', '" + "driver='" + driver + "' WHERE id=1;"
	_, err = db.Exec(quer)

	if err != nil {
		return err
	}
	return nil
}

// DataBaseCheckUserExistance проверяет, внесен ои юзер в базу данных, возврашает его токен если да
func DataBaseCheckUserExistance(login string, password string) (flag int, tknm string, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return 0, "", err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return 0, "", err
	}
	defer db.Close()
	var token string
	if err := db.QueryRow("SELECT token FROM users WHERE lgn = '" + string(login) + "' and psw = '" + string(password) + "';").Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			return 0, "", nil
		}
		return 0, "", err
	}
	return 1, token, nil
}

// DataBasePostUser добавляет запись с новым юзером в базу данных
func DataBasePostUser(login string, password string, token string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO users(lgn, psw, token) VALUES ('` + string(login) + `', '` + string(password) + `', '` + token + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

// DataBaseCheckAuth проверяет, проходит ли юзер аутентификацию
func DataBaseCheckAuth(token string) (flag int, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return 0, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var login string
	if err := db.QueryRow("SELECT lgn FROM users WHERE token = '" + string(token) + "';").Scan(&login); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return 1, nil

}

// DataBasePostUserText заносит в таблицу TextData запись с текстом и токеном юзера, который этот текст отправил
func DataBasePostUserText(token string, text string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO TextData(token, txxt) VALUES ('` + token + `', '` + text + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil

}

// DataBaseGetUserText возвращает все текстовые записи юзера из TextData с конкретным токеном
func DataBaseGetUserText(token string) (answb []flw.TextAnswer, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	quer := "SELECT txxt from TextData where token = '" + token + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		answ := new(flw.TextAnswer)
		err = rows.Scan(&answ.UserText)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		answb = append(answb, *answ)
	}
	return answb, nil

}

// DataBasePostUserBin заносит в таблицу BinData запись с бинарными данными и токеном юзера, который это отправил
func DataBasePostUserBin(token string, bin []byte) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO BinData(token, bin) VALUES ('` + token + `', '` + string(bin) + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil

}

// DataBaseGetUserBin возвращает все бинарные записи юзера из BinData с конкретным токеном
func DataBaseGetUserBin(token string) (answb []flw.BinAnswer, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	quer := "SELECT bin from BinData where token = '" + token + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		answ := new(flw.BinAnswer)
		err = rows.Scan(&answ.UserBin)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		answb = append(answb, *answ)
	}
	return answb, nil

}

// DataBasePostUserCard заносит в таблицу CardData запись с номером карты, именем держателя и cvv с токеном юзера, который это отправил
func DataBasePostUserCard(token string, number string, holder string, cvv string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO CardData(token, CredNumber, CredHolder, CredCVV) VALUES ('` + token + `', '` + number + `', '` + holder + `', '` + cvv + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil

}

// DataBaseGetUserCard возвращает все записи по картам юзера из CardData с конкретным токеном
func DataBaseGetUserCard(token string) (answb []flw.CardAnswer, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	quer := "SELECT credNum, credHolder, credCVV from cardData where token = '" + token + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		answ := new(flw.CardAnswer)
		err = rows.Scan(&answ.UserCardNum, &answ.UserCardName, &answ.UserCardCVV)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		answb = append(answb, *answ)
	}
	return answb, nil

}
