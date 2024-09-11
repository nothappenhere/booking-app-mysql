package main

import (
	"booking-app/database"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var conferenceName string = "Tickify"

const conferenceTicket uint = 100

var currentTicket uint = 100

type users struct {
	FirstName     string
	LastName      string
	Email         string
	Ticket        uint
	Status        string
	PaymentMethod string
	PaymentAt     time.Time
}

var allUser []users
var wg sync.WaitGroup

func main() {
	greeting()

	for {
		firstName, lastName, email, ticket := inputUser()

		isValidFirstName := len(firstName) >= 2
		isValidLastName := len(lastName) >= 2
		isValidEmail := strings.Contains(email, "@")
		isValidTicket := ticket <= currentTicket

		if isValidFirstName && isValidLastName && isValidEmail && isValidTicket {
			currentTicket -= ticket

			user := users{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Ticket:    ticket,
				Status:    "UNPAID",
			}
			fmt.Printf("Thank you for your registration, %s\nyour payment status is %s.\n", user.FirstName, user.Status)

			fmt.Println("~ What would you like to pay with?")
			fmt.Println("1. DANA")
			fmt.Println("2. GOPAY")
			fmt.Println("3. SHOPEEPAY")
			fmt.Println("4. Bank Transfer")

			wg.Add(1)
			go processingPayment(&user)
			wg.Wait()

			wg.Add(1)
			go insertData(&user)
			wg.Wait()

			wg.Add(1)
			go sendTicket(user.Email)
			wg.Wait()

			allUser = append(allUser, user)

			fmt.Printf("Now we have %d tickets remaining.\n", currentTicket)
			fmt.Println("=================================================================")

			wg.Add(1)
			go displayAllUsers()
			wg.Wait()

		} else {
			if !isValidFirstName {
				fmt.Println("First name must be at least 2 characters")
			}
			if !isValidLastName {
				fmt.Println("Last name must be at least 2 characters")
			}
			if !isValidEmail {
				fmt.Println("Email address must contain @ sign")
			}
			if !isValidTicket {
				fmt.Printf("Number of tickets must not exceed the available tickets: %d, you entered: %d\n", currentTicket, ticket)
			}
		}

		if currentTicket == 0 {
			fmt.Println("Thank you for your enthusiasm, but we are out of tickets now.")
			break
		}
	}
}

func greeting() {
	fmt.Println("=================================================================")
	fmt.Printf("Welcome to %s.\n", conferenceName)
	fmt.Printf("Experience the ease of booking your travel and event tickets.\nGetting your tickets is faster, simpler, and secure with %s.\n", conferenceName)
	fmt.Printf("We currently have %d tickets, book your tickets now.\n", conferenceTicket)
	fmt.Println("=================================================================")
}

func inputUser() (string, string, string, uint) {
	fmt.Print("~ First name: ")
	var firstName string
	fmt.Scan(&firstName)

	fmt.Print("~ Last name: ")
	var lastName string
	fmt.Scan(&lastName)

	fmt.Print("~ Email address: ")
	var email string
	fmt.Scan(&email)

	fmt.Print("~ Number of tickets: ")
	var ticket uint
	fmt.Scan(&ticket)

	fmt.Println("=================================================================")
	return firstName, lastName, email, ticket
}

func processingPayment(user *users) {
	defer wg.Done()

	var paymentMethod int
	for {
		fmt.Print("Choose payment method: ")
		fmt.Scan(&paymentMethod)
		fmt.Println("")

		isValidOption := paymentMethod >= 1 && paymentMethod <= 4
		if isValidOption {
			break
		} else {
			fmt.Println("Invalid choice, please choose a valid payment method (1-4).")
		}
	}

	rand.Seed(time.Now().UnixNano())
	var randomNumber int = rand.Intn(90000000) + 10000000
	var formattedPhoneNumber string = fmt.Sprintf("08%d", randomNumber)
	var formattedBankNumber string = fmt.Sprintf("15%d", randomNumber)

	switch paymentMethod {
	case 1:
		fmt.Printf("Here is the number for DANA: %s", formattedPhoneNumber)
	case 2:
		fmt.Println("Here is the number for GOPAY:", formattedPhoneNumber)
	case 3:
		fmt.Println("Here is the number for SHOPEEPAY:", formattedPhoneNumber)
	case 4:
		fmt.Println("Here is the account number for BANK payment:", formattedBankNumber)
	}

	time.Sleep(5 * time.Second)

	user.Status = "PAID"
	user.PaymentMethod = map[int]string{
		1: "DANA",
		2: "GOPAY",
		3: "SHOPEEPAY",
		4: "BANK TRANSFER",
	}[paymentMethod]
	user.PaymentAt = time.Now()

	fmt.Printf("Thank you for your payment, %s.\nYour payment status is %s, payment method is %s at %s\n", user.FirstName, user.Status, user.PaymentMethod, user.PaymentAt.Format("2006-01-02 15:04:05"))
}

func sendTicket(email string) {
	defer wg.Done()

	time.Sleep(2 * time.Second)
	fmt.Printf("We sent your ticket to %s\n", email)
}

func displayAllUsers() {
	defer wg.Done()

	fmt.Println("List of registered users:")
	for _, user := range allUser {
		fmt.Printf("- %s %s, %s, %d ticket(s) - Payment Status: %s\n", user.FirstName, user.LastName, user.Email, user.Ticket, user.Status)
	}
	fmt.Println("=================================================================")
}

func insertData(user *users) {
	defer wg.Done()

	db := database.GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO users(first_name, last_name, email, ticket, status, payment_method, payment_at) VALUES(?, ?, ?, ?, ?, ?, ?)"
	_, err := db.ExecContext(ctx, script, user.FirstName, user.LastName, user.Email, user.Ticket, user.Status, user.PaymentMethod, user.PaymentAt)
	if err != nil {
		fmt.Println("Error inserting data:", err)
		panic(err)
	}
}
