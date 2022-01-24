// https://youtu.be/yyUHQIec83I

package main

import (
	"fmt"
	//"sync"
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

//var wg = sync.WaitGroup{}

func main() {

	greetUsers()

	for remainingTickets > 0 {

		persist.Load("./remainingTickets.tmp", &remainingTickets)
		firstName, lastName, email, userTickets := getUserInput()
		isValidName, isValidEmail, isValidUserTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

		if isValidName && isValidEmail && isValidUserTicketNumber {

			bookTicket(userTickets, firstName, lastName, email)
			//wg.Add(1)
			go sendTicket(userTickets, firstName, lastName, email)

			firstNames := getFirstNames()
			fmt.Printf("The first names of Bookings are: %v\n", firstNames)
		} else if userTickets == remainingTickets {
			fmt.Print("You're about to purchase all the leftover tickets, shame on you (but thanks for giving us money)")
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
	}
	//wg.Wait()
	persist.Save("./remainingTickets.tmp", remainingTickets)
}
func greetUsers() {
	fmt.Printf("welcome to the %v booking cli\n", conferenceName)
	fmt.Printf("We have a total of %v tickets and %v remain available.\n", conferenceTickets, remainingTickets)
}
func getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range Bookings {
		firstNames = append(firstNames, booking.FirstName)
	}
	return firstNames
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
	//wg.Done()
}
