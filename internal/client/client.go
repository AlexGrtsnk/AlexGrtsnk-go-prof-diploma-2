package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	db "github.com/AlexGrtsnk/go-prof-diploma-2/internal/db"
	flw "github.com/AlexGrtsnk/go-prof-diploma-2/internal/json_parser"
)

func ClientTerminalWorking() {
	client := http.Client{}
	fmt.Println("Hello. You are using client-server golang programm. Please choose your option and enter its number:")
	fmt.Println("1. Sign up")
	fmt.Println("2. Sign in")
	apiRunAddr, err := db.DataBaseAPIAddressSelect()
	if err != nil {
		log.Println(err)
		return
	}
	flag := 0
	var token_cookie *http.Cookie
	for {
		var chs int
		var login, password string
		if flag == 0 {
			flag = 1
			_, err := fmt.Scanln(&chs)
			if err != nil {
				flag = 0
				fmt.Println("Input was incorrect, please try again", err)
				continue
			}
			if chs == 1 {
				fmt.Println("Please input your new login:")
				_, err = fmt.Scanln(&login)
				if err != nil {
					flag = 0
					fmt.Println("Login input was incorrect, please try again", err)
					continue
				}
				fmt.Println("Please input your new password:")
				_, err = fmt.Scanln(&password)
				if err != nil {
					flag = 0
					fmt.Println("Pawword input was incorrect, please try again", err)
					continue
				}

				b := new(bytes.Buffer)
				_, err = io.WriteString(b, `{"login":"`+login+`"}, `+`{"password":"`+password+`"}`)
				if err != nil {
					log.Fatal(err)
				}
				request, err := http.NewRequest("POST", apiRunAddr+"/register", b)
				if err != nil {
					log.Fatal(err)
				}
				resp, err := client.Do(request)
				if err != nil {
					fmt.Println("Registration was not sucessfull. Please check your login and password")
					continue
				}
				token_cookie, err = resp.Request.Cookie("exampleCookie")
				if err != nil {
					fmt.Println("Registration was not sucessfull. Please check ypur login and password")
					continue
				}
				continue
			}
			if chs == 2 {
				fmt.Println("Please input your login:")
				_, err = fmt.Scanln(&login)
				if err != nil {
					flag = 0
					fmt.Println("Login input was incorrect, please try again", err)
					continue
				}
				fmt.Println("Please input your password:")
				_, err = fmt.Scanln(&password)
				if err != nil {
					flag = 0
					fmt.Println("Pawword input was incorrect, please try again", err)
					continue
				}

				b := new(bytes.Buffer)
				_, err = io.WriteString(b, `{"login":"`+login+`"}, `+`{"password":"`+password+`"}`)
				if err != nil {
					log.Fatal(err)
				}
				request, err := http.NewRequest("POST", apiRunAddr+"/login", b)
				if err != nil {
					log.Fatal(err)
				}
				resp, err := client.Do(request)
				if err != nil {
					fmt.Println("Entering was not sucessfull. Please check ypur login and password")
					continue
				}
				token_cookie, err = resp.Request.Cookie("exampleCookie")
				if err != nil {
					fmt.Println("Entering was not sucessfull. Please check ypur login and password")
					continue
				}
				continue
			}

		}
		fmt.Println("Your verificastion was sucsefull! Now please choose your option and enter its number: ")
		fmt.Println("1. Add your new text")
		fmt.Println("2. Get all your texts")
		fmt.Println("3. Add your new binary file")
		fmt.Println("4. Get all your binary files")
		fmt.Println("5. Add your new credit card")
		fmt.Println("6. Get information about all your credit cards")
		_, err := fmt.Scanln(&chs)
		if err != nil {
			fmt.Println("Некорректный ввод", err)
			continue
		}
		const (
			postUsersText = iota
			getUsersText
			postUsersBin
			getUsersBin
			postUsersCard
			getUsersCard
		)
		switch chs {
		case postUsersText:
			fmt.Println("Please input your text:")
			var text string
			_, err = fmt.Scanln(&text)
			if err != nil {
				fmt.Println("Text input was incorrect, please try again", err)
				continue
			}
			b := new(bytes.Buffer)
			_, err = io.WriteString(b, `{"user_text":"`+text+`"}`)
			if err != nil {
				log.Fatal(err)
			}
			request, err := http.NewRequest("POST", apiRunAddr+"/text", b)
			request.AddCookie(token_cookie)
			if err != nil {
				log.Fatal(err)
			}
			_, err = client.Do(request)
			if err != nil {
				fmt.Println("Adding text was not sucessfull. Please do it again")
				continue
			}

			continue
		case getUsersText:
			fmt.Println("This is your added texts:")
			b := new(bytes.Buffer)
			request, err := http.NewRequest("GET", apiRunAddr+"/text", b)
			request.AddCookie(token_cookie)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println("Getting text was not sucessfull. Please do it again")
				continue
			}
			dataTmp := new(flw.TextList)
			err = json.NewDecoder(resp.Body).Decode(dataTmp)
			if err != nil {
				fmt.Println("Getting text was not sucessfull. Please do it again")
				continue
			}
			fmt.Println(dataTmp)
			continue
		case postUsersBin:
			fmt.Println("Please input full path to your binary file:")
			var bin string
			_, err = fmt.Scanln(&bin)
			if err != nil {
				fmt.Println("Full path input was incorrect, please try again", err)
				continue
			}
			f, err := os.Open(bin)
			if err != nil {
				panic(err)
			}
			var b1 []byte
			f.Read(b1)
			binTmp := string(b1[:])

			b := new(bytes.Buffer)
			_, err = io.WriteString(b, `{"user_bin":"`+binTmp+`"}`)
			if err != nil {
				log.Fatal(err)
			}
			request, err := http.NewRequest("POST", apiRunAddr+"/bin", b)
			request.AddCookie(token_cookie)
			if err != nil {
				log.Fatal(err)
			}
			_, err = client.Do(request)
			if err != nil {
				fmt.Println("Adding file was not sucessfull. Please do it again")
				continue
			}
			continue
		case getUsersBin:
			fmt.Println("This is your added binaries:")
			b := new(bytes.Buffer)
			request, err := http.NewRequest("GET", apiRunAddr+"/bin", b)
			request.AddCookie(token_cookie)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println("Getting binaries was not sucessfull. Please do it again")
				continue
			}
			dataTmp := new(flw.BinList)
			err = json.NewDecoder(resp.Body).Decode(dataTmp)
			if err != nil {
				fmt.Println("Getting binaries was not sucessfull. Please do it again")
				continue
			}
			fmt.Println(dataTmp)
		case postUsersCard:
			fmt.Println("Please input your card number:")
			var cardNum string
			_, err = fmt.Scanln(&cardNum)
			if err != nil {
				fmt.Println("Card number input was incorrect, please try again", err)
				continue
			}
			fmt.Println("Please input card holder:")
			var cardHolder string
			_, err = fmt.Scanln(&cardHolder)
			if err != nil {
				fmt.Println("Holder input was incorrect, please try again", err)
				continue
			}
			fmt.Println("Please input your card cvv:")
			var cardCVV string
			_, err = fmt.Scanln(&cardCVV)
			if err != nil {
				fmt.Println("Cvv input was incorrect, please try again", err)
				continue
			}
			b := new(bytes.Buffer)
			_, err = io.WriteString(b, `{"user_card_num":"`+cardNum+`"}, `+`{"user_card_name":"`+cardHolder+`"}, `+`{"user_card_cvv":"`+cardCVV+`"}`)
			if err != nil {
				log.Fatal(err)
			}
			request, err := http.NewRequest("POST", apiRunAddr+"/card", b)
			request.AddCookie(token_cookie)
			if err != nil {
				log.Fatal(err)
			}
			_, err = client.Do(request)
			if err != nil {
				fmt.Println("Adding text was not sucessfull. Please do it again")
				continue
			}
			continue
		case getUsersCard:
			fmt.Println("This is your added cards:")
			b := new(bytes.Buffer)
			request, err := http.NewRequest("GET", apiRunAddr+"/card", b)
			request.AddCookie(token_cookie)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println("Getting cards was not sucessfull. Please do it again")
				continue
			}
			dataTmp := new(flw.CardList)
			err = json.NewDecoder(resp.Body).Decode(dataTmp)
			if err != nil {
				fmt.Println("Getting cards was not sucessfull. Please do it again")
				continue
			}
			fmt.Println(dataTmp)
		}
	}
}
