package database

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestInsertData(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO users(first_name, last_name, email, ticket, status, payment_method, payment_at) VALUES(?, ?, ?, ?, ?, ?, ?);"
	stmt, err := db.PrepareContext(ctx, script)
	if err != nil {
		t.Fatalf("Failed to insert data into tabel users: %v", err)
	}

	for i := 0; i < 10; i++ {
		first_name := "John" + strconv.Itoa(i)
		last_name := "Doe" + strconv.Itoa(i)
		email := "john.doe@example.com" + strconv.Itoa(i)
		ticket := 2 + i
		status := "UNPAID"
		payment_method := "Credit Card"
		payment_at := time.Now()

		result, err := stmt.ExecContext(ctx, first_name, last_name, email, ticket, status, payment_method, payment_at)
		if err != nil {
			t.Fatalf("Failed to insert data into tabel users: %v", err)
		}

		lastInsertedID, _ := result.LastInsertId()
		fmt.Println(lastInsertedID)
	}

	fmt.Println("Success insert into table users")
}

func TestQueryData(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, first_name, last_name, email, ticket, status, payment_method, payment_at FROM users"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		t.Fatalf("Failed to query data from table users: %v", err)
	}

	for rows.Next() {
		var id int
		var first_name string
		var last_name string
		var email string
		var ticket int
		var status string
		var payment_method string
		var payment_at time.Time

		err := rows.Scan(&id, &first_name, &last_name, &email, &ticket, &status, &payment_method, &payment_at)
		if err != nil {
			t.Errorf("Failed displaying data from table users: %v", err)
		}

		fmt.Println(id, first_name, last_name, email, ticket, status, payment_method, payment_at)
	}

	defer rows.Close()
}
