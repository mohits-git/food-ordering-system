package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	apiclient "github.com/mohits-git/food-ordering-system/cmd/cli/api_client"
	"github.com/mohits-git/food-ordering-system/cmd/cli/handlers"
)

func main() {
	apiClient := apiclient.NewAPIClient("http://localhost:8080")
	handler := handlers.NewHandlers(apiClient)

	startRepl(handler)
}

func startRepl(handler *handlers.Handlers) {
	jwtToken := ""
	for {
		clearScreen()
		// welcome message
		fmt.Println("Welcome to the Food Ordering System!")
		printActions()
		// Read input
		action := -1
		fmt.Println("Choose an action:")
		fmt.Scan(&action)
		fmt.Println()

		// Process input
		switch action {
		case 0:
			fmt.Println("Exiting...")
			return
		case 1:
			jwtToken = handler.HandleLogin()
		case 2:
			handler.HandleLogout(jwtToken)
		case 3:
			handler.HandleRegisterCustomer()
		case 4:
			handler.HandleViewRestaurants()
		case 5:
			handler.HandleViewRestaurantMenuItems()
		case 6:
			handler.HandleCreateOrder(jwtToken)
		case 7:
			handler.HandleAddMenuItemToOrder(jwtToken)
		case 8:
			handler.HandlePlaceOrderAndGetBill(jwtToken)
		case 9:
			handler.HandlePayBill(jwtToken)
		case 10:
			handler.HandleGetInvoiceById(jwtToken)
		case 11:
			handler.HandleRegisterRestaurantOwner()
		case 12:
			handler.HandleAddRestaurant(jwtToken)
		case 13:
			handler.HandleAddMenuItemToRestaurant(jwtToken)
		case 14:
			handler.HandleUpdateMenuItemAvailability(jwtToken)
		}

		// clear screen
		fmt.Println()
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
		clearScreen()
	}
}

func printActions() {
	menu := `

  Available actions:
  0. Exit
  1. Login
  2. Logout

  Customer Actions:
  3. Register as Customer
  4. View Restaurants
  5. View Restaurants Menu Items
  6. Create Order
  7. Add Menu Item to Order
  8. Place Order and Get Bill
  9. Pay for Bill
  10. Get Bill by Id

  Restaurant Owner Actions:
  11. Register as Restaurant Owner
  12. Add Restaurant
  13. Add Menu Item to Restaurant
  14. Update Menu Item Availability

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
