package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type Item struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"_"`
	UpdatedAt    time.Time `json:"_"`
	Title        string    `json:"title"`
	Notes        string    `json:"notes"`
	SellerID     uuid.UUID `json:"seller"`
	PriceInCents int64     `json:"price_in_cents"`
}

func (i *Item) Create(conn *pgx.Conn, userID string) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty.")
	}
	if i.PriceInCents < 0 {
		i.PriceInCents = 0
	}
	now := time.Now()

	row := conn.QueryRow(context.Background(), "INSERT INTO item (title, notes, seller_id, price_in_cents, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, seller_id", i.Title, i.Notes, userID, i.PriceInCents, now, now)

	err := row.Scan(&i.ID, &i.SellerID)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the item")
	}

	return nil
}

func GetAllItems(conn *pgx.Conn) ([]Item, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, title, notes, seller_id, price_in_cents FROM item")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error getting items")
	}

	var items []Item
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ID, &item.Title, &item.Notes, &item.SellerID, &item.PriceInCents)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func GetItemsBeingSoldByUser(userID string, conn *pgx.Conn) ([]Item, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, title, price_in_cents, notes, seller_id FROM item WHERE seller_id = $1", userID)
	if err != nil {
		fmt.Printf("Error getting items %v", err)
		return nil, fmt.Errorf("There was an error getting the items")
	}

	var items []Item
	for rows.Next() {
		i := Item{}
		err = rows.Scan(&i.ID, &i.Title, &i.PriceInCents, &i.Notes, &i.SellerID)
		if err != nil {
			fmt.Printf("Error scaning item: %v", err)
			continue
		}
		items = append(items, i)
	}

	return items, nil
}

func (i *Item) Update(conn *pgx.Conn) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty")
	}
	if i.PriceInCents < 0 {
		i.PriceInCents = 0
	}
	now := time.Now()
	_, err := conn.Exec(context.Background(), "UPDATE item SET title=$1, notes=$2, price_in_cents=$3, updated_at=$4 WHERE id=$5", i.Title, i.Notes, i.PriceInCents, now, i.ID)

	if err != nil {
		fmt.Printf("Error updating item: (%v)", err)
		return fmt.Errorf("Error updating item")
	}

	return nil
}

func FindItemById(id uuid.UUID, conn *pgx.Conn) (Item, error) {
	row := conn.QueryRow(context.Background(), "SELECT title, notes, seller_id, price_in_cents FROM item WHERE id=$1", id)
	item := Item{
		ID: id,
	}
	err := row.Scan(&item.Title, &item.Notes, &item.SellerID, &item.PriceInCents)
	if err != nil {
		return item, fmt.Errorf("The item doesn't exist")
	}

	return item, nil
}
