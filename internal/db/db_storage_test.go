package databasestorage

import (
	"log"
	"testing"

	apcfg "github.com/AlexGrtsnk/go-prof-diploma-2/internal/app_config"

	"github.com/caarlos0/env"
)

func TestNewDB(t *testing.T) {
	_, err := NewDB()
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestRunMigrateScripts(t *testing.T) {
	err := RunMigrateScripts(nil)
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBaseCheckAuth(t *testing.T) {
	_, err := DataBaseCheckAuth("")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}
func TestDataBaseCheckUserExistance(t *testing.T) {
	_, _, err := DataBaseCheckUserExistance("", "")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}
func TestDataBasePostUser(t *testing.T) {
	err := DataBasePostUser("", "", "")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBaseCfg(t *testing.T) {
	err := DataBaseCfg("", "", "", "")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBasePingHandler(t *testing.T) {
	err := DataBasePingHandler()
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBasePing(t *testing.T) {
	err := DataBasePing("", "")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBasePostUserText(t *testing.T) {
	err := DataBasePostUserText("", "")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBaseGetUserText(t *testing.T) {
	_, err := DataBaseGetUserText("")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBasePostUserBin(t *testing.T) {
	err := DataBasePostUserBin("", nil)
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}
func TestDataBaseGetUserBin(t *testing.T) {
	_, err := DataBaseGetUserBin("")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBasePostUserCard(t *testing.T) {
	err := DataBasePostUserCard("", "", "", "")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBaseGetUserCard(t *testing.T) {
	_, err := DataBaseGetUserCard("")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBaseSelfConfigGet(t *testing.T) {
	_, _, err := DataBaseSelfConfigGet()
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBaseSelfConfigUpdate(t *testing.T) {
	_, _, err := DataBaseSelfConfigGet()
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestNewDBGood(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did  panic")
		}
	}()
	_ = DataBaseStartConfig(":8080")
	_, err := NewDB()
	if err != nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestRunMigrateScriptsBad(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_ = DataBaseStartConfig(":8080")
	err := RunMigrateScripts(nil)
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestDataBaseCreateShortURLPageCfgGood(t *testing.T) {
	var cfg apcfg.Config
	err := env.Parse(&cfg)
	flagRunAddr, apiRunAddr, fileName, databaseDSN, enableHTTPS, config, trustedSubnet := apcfg.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ServerAddress != "" {
		flagRunAddr = "8080"
	}
	if cfg.BaseURL != "" {
		apiRunAddr = cfg.BaseURL
	}
	if cfg.FileStoragePath != "" {
		fileName = cfg.FileStoragePath
	}
	if cfg.DatabaseDSN != "" {
		databaseDSN = cfg.DatabaseDSN
	}
	if cfg.EnableHTTPS != true {
		enableHTTPS = cfg.EnableHTTPS
	}
	log.Println(enableHTTPS)
	log.Println(config)
	log.Println(trustedSubnet)
	err = DataBaseStartConfig(databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if databaseDSN != "localhost" {
		err = DataBasePingHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
	err = DataBaseCfg(flagRunAddr, apiRunAddr, fileName, trustedSubnet)
	if err == nil {
		log.Fatal(err)
	}

	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}
