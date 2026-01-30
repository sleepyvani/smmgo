# SMM (SOCIAL MEDIA MARKETING) - Go

A lightweight, zero-dependency Go client for connecting to compatible SMM (Social Media Marketing) Panels.

## Features
- **Zero Dependencies**: Uses native `net/http` package.
- **Go Support**: Works with standard Go versions.

## Installation

```bash
go get github.com/sleepyvani/smmgo
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/sleepyvani/smmgo"
)

func main() {
	client := smm.New("https://smm-provider.com/api/v2", "YOUR_API_KEY")

	// Get Balance
	balance, err := client.GetBalance()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Balance:", balance)

	// Add Order
	/*
	order, err := client.AddOrder(smm.AddOrderParams{
		Service:  123,
		Link:     "https://example.com",
		Quantity: 1000,
	})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Order:", order)
	}
	*/
}
```

## API Reference
- `GetServices()`
- `AddOrder(params AddOrderParams)`
- `GetStatus(orderID interface{})`
- `GetMultiStatus(orderIDs []interface{})`
- `CreateRefill(orderID interface{})`
- `GetRefillStatus(refillID interface{})`
- `CancelOrders(orderIDs []interface{})`
- `GetBalance()`

## License
MIT
