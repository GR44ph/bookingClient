// https://youtu.be/yyUHQIec83I

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/GR44ph/bookingClient/helper"
	"github.com/GR44ph/bookingClient/persist"
)

var conferenceName = "Go Conference"

const conferenceTickets uint = 50

var remainingTickets uint = conferenceTickets
var Bookings = make([]UserData, 0)

type UserData struct {
	FirstName   string
	LastName    string
	Email       string
	UserTickets uint
}

var wg = sync.WaitGroup{}

func main() {
	switch os.Args[1] {
	default:
		fmt.Println("Please use one of these command-arguments: order, reset, list")
	case "order":
		order()
	case "reset":
		reset()
	case "list":
		list()
	}
}
func order() {

	persist.Load("./remainingTickets.tmp", &remainingTickets)
	greetUsers()
	firstName, lastName, email, userTickets := getUserInput()
	isValidName, isValidEmail, isValidUserTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

	if isValidName && isValidEmail && isValidUserTicketNumber {

		bookTicket(userTickets, firstName, lastName, email)
		wg.Add(1)
		go sendTicket(userTickets, firstName, lastName, email)
	} else {
		if !isValidName {
			fmt.Println("The first name or last name you've entered is too short (at least 2 characters)")
		}
		if !isValidEmail {
			fmt.Println("The Email Address you've entered does not contain an @ sign")
		}
		if !isValidUserTicketNumber {
			fmt.Println("The Number of tickets you've entered is invalid")
		}
	}
	wg.Wait()
	persist.Save("./remainingTickets.tmp", remainingTickets)
}
func greetUsers() {
	fmt.Printf("welcome to the %v booking cli\n", conferenceName)
	fmt.Printf("We have a total of %v tickets and %v remain available.\n", conferenceTickets, remainingTickets)
}
func getUserInput() (string, string, string, uint) {
	var firstName string
	var lastName string
	var email string
	var userTickets uint
	fmt.Print("Enter your first name: ")
	fmt.Scan(&firstName)
	fmt.Print("Enter your last name: ")
	fmt.Scan(&lastName)
	fmt.Print("Enter your email address: ")
	fmt.Scan(&email)
	fmt.Print("Enter number of tickets: ")
	fmt.Scan(&userTickets)
	return firstName, lastName, email, userTickets
}

func reset() {
	fmt.Println("Resetting...")
	os.Remove("bookings.tmp")
	fmt.Println("Deleted bookings.tmp")
	os.Remove("remainingTickets.tmp")
	fmt.Println("Deleted remainingTickets.tmp")
	fmt.Println("Done!")
}

func list() {
	persist.Load("./bookings.tmp", &Bookings)
	jsonF, _ := json.Marshal(Bookings)
	fmt.Println(string(jsonF))
}
func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets -= userTickets

	var userData = UserData{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		UserTickets: userTickets,
	}
	persist.Load("./bookings.tmp", &Bookings)
	Bookings = append(Bookings, userData)
	persist.Save("./bookings.tmp", Bookings)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, conferenceName)
}
func sendTicket(userTickets uint, firstName string, lastName string, email string) {
	time.Sleep(10 * time.Second)
	var ticket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	fmt.Println("##################")
	fmt.Printf("Sending ticket:\n %v to email address %v\n", ticket, email)
	fmt.Println("##################")
	wg.Done()
}
