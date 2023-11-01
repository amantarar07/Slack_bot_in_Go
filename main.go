package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/shomali11/slacker"
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

func main() {


//xoxe.xoxp-1-Mi0yLTYxMjAyNjMxOTMyMjItNjE0OTg5ODU1Nzc3Ni02MTI3NzUxMzc5OTIyLTYxNTA3NDUwMjI2NzItNGQzMmU4OTY1MDJlNDFlZWZmM2Q2MmQ3NDA1YTg4ZTMxMWU5MGUxMTNjMWViY2MxNTgzNDVhZjI0NDYwYzRmZA


//xapp-1-A0643MS8HS5-6124193798709-34ccae9593ab95f8feddd4812771c4db712fdae42a8aa1d69b3de004624ee5c7


option:=slacker.ClientOption(func(cd *slacker.ClientDefaults) {})

	bot := slacker.NewClient("xoxb-6120263193222-6139809551649-JELwbo7VfnwJimyzUM30GESC", "xoxe.xoxp-1-Mi0yLTYxMjAyNjMxOTMyMjItNjE0OTg5ODU1Nzc3Ni02MTI3NzUxMzc5OTIyLTYxNTA3NDUwMjI2NzItNGQzMmU4OTY1MDJlNDFlZWZmM2Q2MmQ3NDA1YTg4ZTMxMWU5MGUxMTNjMWViY2MxNTgzNDVhZjI0NDYwYzRmZA")
	go printCommandEvents(bot.CommandEvents())
	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			w.Reply("pongy pongy!!")
		},
	})


	bot.Command("user-info", &slacker.CommandDefinition{Handler:func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
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
	},})


	bot.Command("hello", &slacker.CommandDefinition{
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {

			user, err:= bot.APIClient().GetUserByEmail(bc.Event().UserID)
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

		log.Fatal(err)
	}
}
