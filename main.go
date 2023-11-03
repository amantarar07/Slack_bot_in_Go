package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {

	for event := range analyticsChannel {

		fmt.Println("Command events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)

	}

}
var wolframClient *wolfram.Client


func main() {

	//xoxe.xoxp-1-Mi0yLTYxMjAyNjMxOTMyMjItNjE0OTg5ODU1Nzc3Ni02MTI3NzUxMzc5OTIyLTYxNTA3NDUwMjI2NzItNGQzMmU4OTY1MDJlNDFlZWZmM2Q2MmQ3NDA1YTg4ZTMxMWU5MGUxMTNjMWViY2MxNTgzNDVhZjI0NDYwYzRmZA

	//xapp-1-A0643MS8HS5-6124193798709-34ccae9593ab95f8feddd4812771c4db712fdae42a8aa1d69b3de004624ee5c7

	// option:=slacker.ClientOption(func(cd *slacker.ClientDefaults) {})
	//xapp-1-A064GGLLEV7-6125397036103-ef62d47c4ca4e25ee8926296f10342bd094f288e3de9f7fa562aa4b8ed9b1174

	godotenv.Load(".env")
	fmt.Println("envs data", os.Getenv("SLACK_BOT_TOKEN"))

	bot := slacker.NewClient("xoxb-6120263193222-6133326496150-8owcAgQFjALikQRyYoiCOqMk", "xapp-1-A064GGLLEV7-6125397036103-ef62d47c4ca4e25ee8926296f10342bd094f288e3de9f7fa562aa4b8ed9b1174")
	client := witai.NewClient(os.Getenv("WIT_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

	go printCommandEvents(bot.CommandEvents())

	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			w.Reply("pongy pongy!!")
		},
	})

	bot.Command("btao- <message>", &slacker.CommandDefinition{

		Description: "send any question to wolfram!",
		Examples:    []string{"who is the president of india"},
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			query := r.Param("message")
			fmt.Println("query", query)
			message, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})

			data, _ := json.MarshalIndent(message, "", "    ")
			rough := string(data[:])

			fmt.Println("rough: ", rough)
			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			fmt.Println("value", value)
			fmt.Println("message", message)
			ans := value.String()
			res,err:=wolframClient.GetSpokentAnswerQuery(ans,wolfram.Metric,1000)
			if err!=nil{

				fmt.Println("there is an error")
			}

			fmt.Println("value",value)


			w.Reply(res)
		},
	})

	bot.Command("user-info", &slacker.CommandDefinition{Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
		// Get the user ID from the command arguments (e.g., "@username")
		// userID := r.Param("user-email")

		// Use the GetUserInfo method to retrieve user information
		user, err := bot.APIClient().GetUserByEmail(bc.Event().UserID)
		if err != nil {
			w.Reply("User not found or error occurred.")
			return
		}

		w.Reply("User ID: " + user.ID)
		w.Reply("Username: " + user.Name)
		w.Reply("Full Name: " + user.RealName)
	}})

	bot.Command("hello", &slacker.CommandDefinition{
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {

			user, err := bot.APIClient().GetUserByEmail(bc.Event().UserID)
			if err != nil {
				fmt.Println("error in getting user", err)
				return
			}

			w.Reply("Hello, <@" + user.Name + ">! How can I assist you today?")
		},
	})

	bot.Command("help", &slacker.CommandDefinition{
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			// List of available commands
			commands := []string{"ping", "hello", "help"}
			w.Reply("Available commands: " + strings.Join(commands, ", "))
		},
	})

	bot.Command("echo <text>", &slacker.CommandDefinition{
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			text := r.Param("text")
			w.Reply("You said: " + text)
		},
	})

	bot.Command("quote", &slacker.CommandDefinition{
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			// Generate a random quote or fetch it from an external source
			quote := "Bhuperder Yogi!!"
			w.Reply(quote)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {

		fmt.Println("error in listening", err)
		log.Fatal(err.Error())
	}
}
