package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	apcfg "github.com/AlexGrtsnk/go-prof-diploma-2/internal/app_config"
	ath "github.com/AlexGrtsnk/go-prof-diploma-2/internal/authentication"

	cks "github.com/AlexGrtsnk/go-prof-diploma-2/internal/cookies"
	db "github.com/AlexGrtsnk/go-prof-diploma-2/internal/db"
	gzp "github.com/AlexGrtsnk/go-prof-diploma-2/internal/gzp"
	flw "github.com/AlexGrtsnk/go-prof-diploma-2/internal/json_parser"
	lg "github.com/AlexGrtsnk/go-prof-diploma-2/internal/logger"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// registrateNewUserPage страница, по адресу которой происходит регистрация пользователя
func registrateNewUserPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		reader, err := gzp.GzipFormatHandlerJSON(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		var tokenTmp string
		_, err = cks.GetCookieHandler(w, r)
		if err != nil {
			token, err := ath.BuildJWTString()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			tokenTmp = token
		}
		var ath flw.Auth
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &ath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flag, _, err := db.DataBaseCheckUserExistance(ath.Login, ath.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 1 {
			w.WriteHeader(http.StatusConflict)
		} else {
			err = db.DataBasePostUser(ath.Login, ath.Password, tokenTmp)
			qwe := cks.SetCookieHandler(w, r, tokenTmp)
			http.SetCookie(w, qwe)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			w.WriteHeader(http.StatusOK)
		}

	}
}

// authentificateUserPage страница, по адресу которой происходит аутентификация пользователя
func authentificateUserPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		reader, err := gzp.GzipFormatHandlerJSON(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var ath flw.Auth
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &ath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flag, token, err := db.DataBaseCheckUserExistance(ath.Login, ath.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		cks.SetCookieHandler(w, r, token)
		w.WriteHeader(http.StatusOK)
	}

}

// uploadNewTextPage страница, по адресу которой при методе POST пользователь заносит текстовые данные на сервер, при GET получает свои записи
func uploadNewTextPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		reader, err := gzp.GzipFormatHandlerJSON(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var txxt flw.TextAnswer
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &txxt); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = db.DataBasePostUserText(token, txxt.UserText)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if r.Method == http.MethodGet {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		tempText, err := db.DataBaseGetUserText(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if err = json.NewEncoder(w).Encode(tempText); err != nil {
			log.Panic(err)
		}

	}

}

// uploadNewBinPage страница, по адресу которой при методе POST пользователь заносит бинарные данные на сервер, при GET получает свои записи
func uploadNewBinPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		reader, err := gzp.GzipFormatHandlerJSON(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var biin flw.BinAnswer
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &biin); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = db.DataBasePostUserBin(token, biin.UserBin)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if r.Method == http.MethodGet {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		tempBin, err := db.DataBaseGetUserBin(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if err = json.NewEncoder(w).Encode(tempBin); err != nil {
			log.Panic(err)
		}
	}

}

// uploadNewBinPage страница, по адресу которой при методе POST пользователь заносит текстовые данные на сервер, при GET получает свои записи
func uploadNewCreditCardPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		reader, err := gzp.GzipFormatHandlerJSON(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var cardData flw.CardAnswer
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &cardData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = db.DataBasePostUserCard(token, cardData.UserCardNum, cardData.UserCardName, cardData.UserCardCVV)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if r.Method == http.MethodGet {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		tempCard, err := db.DataBaseGetUserCard(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if err = json.NewEncoder(w).Encode(tempCard); err != nil {
			log.Panic(err)
		}
	}

}

// Run определяет необходимые для работы приложения системные переменные, настраивает базу данных и запускает сам сервер
func Run() (*http.Server, bool) {
	var cfg apcfg.Config
	var srv = http.Server{}
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
	if !cfg.EnableHTTPS {
		enableHTTPS = cfg.EnableHTTPS
	}
	if cfg.Config != "" {
		config = cfg.Config
	}
	if cfg.TrustedSubnet != "" {
		trustedSubnet = cfg.TrustedSubnet
	}
	log.Println(cfg)
	configFile, err := os.Open(config)
	if err != nil {
		fmt.Println("no file was found")
	}
	var setting flw.Setting
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&setting); err != nil {
		fmt.Println("wrong data in file. using old types")
	}
	if flagRunAddr == "8080" && setting.ServerAddress != "" {
		flagRunAddr = setting.ServerAddress
	}
	if apiRunAddr == "http://localhost:8080" && setting.BaseURL != "" {
		apiRunAddr = setting.BaseURL
	}
	if fileName == "text.txt" && setting.FileStoragePath != "" {
		fileName = setting.FileStoragePath
	}
	if databaseDSN == "localhost" && setting.DataBaseDSN != "" {
		databaseDSN = setting.DataBaseDSN
	}
	if !enableHTTPS && !setting.EnableHTTPS {
		enableHTTPS = setting.EnableHTTPS
	}
	if trustedSubnet == "ns" && setting.TrustedSubnet != "" {
		trustedSubnet = setting.TrustedSubnet
	}
	err = db.DataBaseStartConfig(databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if databaseDSN != "localhost" {
		err = db.DataBasePingHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
	err = db.DataBaseCfg(flagRunAddr, apiRunAddr, fileName, trustedSubnet)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("where postgres is hosted:", databaseDSN)
	log.Println("where db is held", fileName)
	log.Println("Running server on", flagRunAddr)
	log.Println("Running api on", apiRunAddr)
	mux1 := mux.NewRouter()
	mux1.HandleFunc(`/api/user/register`, lg.WithLogging(registrationHandler()))
	mux1.HandleFunc(`/api/user/login`, lg.WithLogging(loginHandler()))
	mux1.HandleFunc(`/api/user/text`, lg.WithLogging(TextUploadHandler()))
	mux1.HandleFunc(`/api/user/bin`, lg.WithLogging(BinUploadHandler()))
	mux1.HandleFunc(`/api/user/card`, lg.WithLogging(CardUploadHandler()))
	mux1.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux1.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux1.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux1.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux1.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	mux1.Handle("/debug/pprof/{cmd}", http.HandlerFunc(pprof.Index))
	srv.Addr = flagRunAddr
	srv.Handler = mux1
	return &srv, enableHTTPS
}

func registrationHandler() http.Handler {
	fn := registrateNewUserPage
	return http.HandlerFunc(fn)
}

func loginHandler() http.Handler {
	fn := authentificateUserPage
	return http.HandlerFunc(fn)
}

func TextUploadHandler() http.Handler {
	fn := uploadNewTextPage
	return http.HandlerFunc(fn)
}

func BinUploadHandler() http.Handler {
	fn := uploadNewBinPage
	return http.HandlerFunc(fn)
}

func CardUploadHandler() http.Handler {
	fn := uploadNewCreditCardPage
	return http.HandlerFunc(fn)
}
