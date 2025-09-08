package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	apiclient "github.com/mohits-git/food-ordering-system/cmd/cli/api_client"
	"github.com/mohits-git/food-ordering-system/cmd/cli/handlers"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

var jwtToken string
var userClaims authctx.UserClaims

func main() {
	apiClient := apiclient.NewAPIClient("http://localhost:8080")
	handler := handlers.NewHandlers(apiClient)

	startRepl(handler)
}

func startRepl(handler *handlers.Handlers) {
	clearScreen()

	for {
		if jwtToken == "" {
			whenUserNotLoggedIn(handler)
		} else if userClaims.Role == "customer" {
			whenCustomerLoggedIn(handler)
		} else if userClaims.Role == "owner" {
			whenRestaurantOwnerLoggedIn(handler)
		} else {
			fmt.Println("Unknown user role. Logging out for safety.")
			jwtToken = ""
			userClaims = authctx.UserClaims{}
		}
		fmt.Printf("\nPress Enter to continue...\n")
		fmt.Scanln()
		clearScreen()
	}
}

func whenUserNotLoggedIn(handlers *handlers.Handlers) {
	printLoggedOutMenu()

	action := -1
	fmt.Println("Choose an action:")
	fmt.Scan(&action)
	fmt.Println()

	clearScreen()
	switch action {
	case 0:
		fmt.Println("Exiting...")
		os.Exit(0)
	case 1:
		jwtToken, userClaims = handlers.HandleLogin()
	case 2:
		handlers.HandleRegisterCustomer()
	case 3:
		handlers.HandleRegisterRestaurantOwner()
	}
}

func printLoggedOutMenu() {
	menu := `
  Welcome to the Food Ordering System!

  Available actions:
  0. Exit
  1. Login
  2. Register as Customer
  3. Register as Restaurant Owner
 
`
	fmt.Println(menu)
}

func whenCustomerLoggedIn(handlers *handlers.Handlers) {
	printCustomerMenu()

	action := -1
	fmt.Println("Choose an action:")
	fmt.Scan(&action)
	fmt.Println()

	clearScreen()
	switch action {
	case 0:
		fmt.Println("Exiting...")
		os.Exit(0)
	case 1:
		handlers.HandleViewRestaurants()
	case 2:
		handlers.HandleViewRestaurantMenuItems()
	case 3:
		handlers.HandlePlaceOrder(jwtToken)
	case 4:
		handlers.HandleLogout(jwtToken)
		jwtToken = ""
		userClaims = authctx.UserClaims{}
	}

	// case 4:
	// 	handlers.HandleAddMenuItemToOrder(jwtToken)
	// case 5:
	// 	handlers.HandlePlaceOrderAndGetBill(jwtToken)
	// case 6:
	// 	handlers.HandlePayBill(jwtToken)
	// case 7:
	// 	handlers.HandleGetInvoiceById(jwtToken)
}

func printCustomerMenu() {
	menu := `
  Welcome Customer!
  
  Available actions:
  0. Exit
  1. View Restaurants
  2. View Restaurants Menu Items
  3. Place Order
  4. Logout
 
`
	// 4. Add Menu Item to Order
	// 5. Place Order and Get Bill
	// 6. Pay for Bill
	// 7. Get Bill by Id

	fmt.Println(menu)
}

func whenRestaurantOwnerLoggedIn(handlers *handlers.Handlers) {
	printRestaurantOwnerMenu()

	action := -1
	fmt.Println("Choose an action:")
	fmt.Scan(&action)
	fmt.Println()

	clearScreen()
	switch action {
	case 0:
		fmt.Println("Exiting...")
		os.Exit(0)
	case 1:
		handlers.HandleAddRestaurant(jwtToken)
	case 2:
		handlers.HandleAddMenuItemToRestaurant(jwtToken)
	case 3:
		handlers.HandleUpdateMenuItemAvailability(jwtToken)
	case 4:
		handlers.HandleLogout(jwtToken)
		jwtToken = ""
		userClaims = authctx.UserClaims{}
	}
}

func printRestaurantOwnerMenu() {
	menu := `
  Welcome Restaurant Owner!
 
  Available actions:
  0. Exit
  1. Add Restaurant
  2. Add Menu Item to Restaurant
  3. Update Menu Item Availability
  4. Logout
 
`
	fmt.Println(menu)
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
