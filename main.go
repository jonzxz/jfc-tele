package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	md "github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"github.com/joho/godotenv"
	jfcModels "github.com/jonzxz/jfc/models"

	tele "gopkg.in/telebot.v3"
)

var (
	menu     = &tele.ReplyMarkup{ResizeKeyboard: true}
	selector = &tele.ReplyMarkup{ResizeKeyboard: true}
	btnHelp  = menu.Text("help")
	//btnSettings    = menu.Text("Settings")
	btnCheckPeople = menu.Text("Check people")

	btnPrev = selector.Data("<-", "prev")
	btnNext = selector.Data("->", "next")
)

func loadTelegramApiKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("TELEGRAM_API_KEY")
}

func main() {

	pref := tele.Settings{
		Token:  loadTelegramApiKey(),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	initBtns()
	bot, err := tele.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/menu", func(c tele.Context) error {
		return c.Send("Displaying menu..", menu)
	})

	bot.Handle(&btnCheckPeople, func(c tele.Context) error {
		response := getJfcPeople()

		returnMsg := formatGetPeople(response)

		if err != nil {
			log.Fatal(err)
		}

		// required because of table
		return c.Send(returnMsg, &tele.SendOptions{
			ParseMode: tele.ModeHTML,
		})
	})

	bot.Handle(&btnHelp, func(c tele.Context) error {
		return c.Edit("Here is some help: ...")
	})

	bot.Handle(&btnPrev, func(c tele.Context) error {
		return c.Respond()
	})

	bot.Start()
}

func initBtns() {

	menu.Reply(
		menu.Row(btnCheckPeople),
		menu.Row(btnHelp),
		//menu.Row(btnSettings),
	)

	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)
}

// Returns Name, TGID, Household in a HTML markdown table
func formatGetPeople(people []jfcModels.Person) string {

	allPeopleData := [][]string{}

	for _, p := range people {
		personData := []string{}
		personData = append(personData, p.Name)
		personData = append(personData, p.TelegramId)
		personData = append(personData, p.Household)

		allPeopleData = append(allPeopleData, personData)
	}

	formattedData, err := md.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("Name", "TelegramId", "Household").
		Format(allPeopleData)

	if err != nil {
		log.Fatal(err)
	}

	return "<pre>" + formattedData + "</pre>"
}

func getJfcPeople() []jfcModels.Person {

	var people []jfcModels.Person

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/people/list", nil)
	if err != nil {
		fmt.Println("new request err")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(responseBytes, &people)

	if err != nil {
		fmt.Println(err)
	}

	return people
}

//func sendToJfcAddPayment() {
//jsonStruct := jfcModels.Payment{
//Type:        "conservancy",
//Remarks:     "Test from go",
//TotalAmount: 123.45,
//Household:   "678",
//}

//payloadBytes, err := json.Marshal(jsonStruct)
//if err != nil {
//// handle err
//}

//body := bytes.NewReader(payloadBytes)

//req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/payments/add", body)
//if err != nil {
//fmt.Println("new request err")
//}
//req.Header.Set("Content-Type", "application/json")

//resp, err := http.DefaultClient.Do(req)
//if err != nil {
//fmt.Println(err)
//}
//defer resp.Body.Close()

//}
