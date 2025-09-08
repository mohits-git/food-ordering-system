package handlers

import (
	"fmt"

	apiclient "github.com/mohits-git/food-ordering-system/cmd/cli/api_client"
	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type Handlers struct {
	apiClient *apiclient.APIClient
}

func NewHandlers(apiClient *apiclient.APIClient) *Handlers {
	return &Handlers{
		apiClient: apiClient,
	}
}

func (h *Handlers) HandleLogin() string {
	var email string
	var password string

	// enter email
	fmt.Println("Enter email:")
	fmt.Scanln(&email)

	// enter password
	fmt.Println("Enter password:")
	fmt.Scanln(&password)

	token, err := h.apiClient.PostLogin(email, password)
	if err != nil {
		fmt.Println("Error while login:", err)
		return ""
	}

	fmt.Printf("Login successful.\nToken: %s\n", token)
	return token
}

func (h *Handlers) HandleLogout(token string) {
	err := h.apiClient.PostLogout(token)
	if err != nil {
		fmt.Println("Error while logout:", err)
		return
	}

	fmt.Println("Logout successful.")
}

func (h *Handlers) handleCreateUser(role string) {
	var name, email, password string

	// enter name
	fmt.Println("Enter name:")
	fmt.Scanln(&name)

	// enter email
	fmt.Println("Enter email:")
	fmt.Scanln(&email)

	// enter password
	fmt.Println("Enter password:")
	fmt.Scanln(&password)

	userID, err := h.apiClient.PostUser(name, email, password, role)
	if err != nil {
		fmt.Println("Error while creating user:", err)
		return
	}

	fmt.Printf("User created successfully with ID: %d\n", userID)
}

func (h *Handlers) HandleRegisterCustomer() {
	h.handleCreateUser("customer")
}

func (h *Handlers) HandleRegisterRestaurantOwner() {
	h.handleCreateUser("owner")
}

func (h *Handlers) HandleViewRestaurants() {
	restaurants, err := h.apiClient.GetRestaurants()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error while fetching restaurants:", err)
		return
	}

	if len(restaurants) == 0 {
		fmt.Println("No restaurants found.")
		return
	}

	fmt.Println("Restaurants:")
	for _, r := range restaurants {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
}

func (h *Handlers) handleViewMenuItemsByRestaurantId(restaurantId int) {
	menuItems, err := h.apiClient.GetMenuItems(restaurantId)
	if err != nil {
		fmt.Println("Error while fetching menu items:", err)
		return
	}

	if len(menuItems) == 0 {
		fmt.Println("No menu items found for this restaurant.")
		return
	}

	fmt.Printf("Menu Items for Restaurant ID %d:\n", restaurantId)
	for _, item := range menuItems {
		availability := "Unavailable"
		if item.Available {
			availability = "Available"
		}
		fmt.Printf("ID: %d, Name: %s, Price: %.2f, Availability: %s\n", item.ID, item.Name, item.Price, availability)
	}
}

func (h *Handlers) HandleViewRestaurantMenuItems() {
	var restaurantId int
	fmt.Println("Enter Restaurant ID:")
	fmt.Scanln(&restaurantId)

	h.handleViewMenuItemsByRestaurantId(restaurantId)
}

func (h *Handlers) HandleAddRestaurant(token string) {
	var name string

	// enter name
	fmt.Println("Enter restaurant name:")
	fmt.Scanln(&name)

	restaurantID, err := h.apiClient.PostRestaurants(name, token)
	if err != nil {
		fmt.Println("Error while adding restaurant:", err)
		return
	}

	fmt.Printf("Restaurant added successfully with ID: %d\n", restaurantID)
}

func (h *Handlers) HandleAddMenuItemToRestaurant(token string) {
	var restaurantId int
	var name string
	var price float64
	var availableInput string
	var available bool

	// enter restaurant ID
	fmt.Println("Enter Restaurant ID:")
	fmt.Scanln(&restaurantId)

	// enter name
	fmt.Println("Enter menu item name:")
	fmt.Scanln(&name)

	// enter price
	fmt.Println("Enter menu item price:")
	fmt.Scanln(&price)
	fmt.Println("Price: ", price)

	// enter availability
	fmt.Println("Is the item available? (yes/no):")
	fmt.Scanln(&availableInput)
	if availableInput == "yes" {
		available = true
	} else {
		available = false
	}

	menuItemID, err := h.apiClient.PostMenuItem(restaurantId, name, price, available, token)
	if err != nil {
		fmt.Println("Error while adding menu item:", err)
		return
	}

	fmt.Printf("Menu item added successfully with ID: %d\n", menuItemID)
}

func (h *Handlers) HandleUpdateMenuItemAvailability(token string) {
	var menuItemId int
	var availableInput string
	var available bool

	// enter menu item ID
	fmt.Println("Enter Menu Item ID:")
	fmt.Scanln(&menuItemId)

	// enter availability
	fmt.Println("Is the item available? (yes/no):")
	fmt.Scanln(&availableInput)
	if availableInput == "yes" {
		available = true
	} else {
		available = false
	}

	err := h.apiClient.PatchMenuItemAvailability(menuItemId, available, token)
	if err != nil {
		fmt.Println("Error while updating menu item availability:", err)
		return
	}

	fmt.Println("Menu item availability updated successfully.")
}

func (h *Handlers) HandleCreateOrder(token string) {
	var restaurantId int

	// enter restaurant ID
	fmt.Println("Enter Restaurant ID:")
	fmt.Scanln(&restaurantId)

	// print restaurant menu items
	h.handleViewMenuItemsByRestaurantId(restaurantId)

	// loop to add menu items to order
	orderItems := []domain.OrderItem{}
	for {
		var menuItemId int
		fmt.Println("Enter Menu Item ID to add to order (0 to finish):")
		fmt.Scanln(&menuItemId)
		if menuItemId == 0 {
			break
		}

		var quantity int
		fmt.Println("Enter quantity:")
		fmt.Scanln(&quantity)
		orderItems = append(orderItems, domain.OrderItem{
			MenuItemID: menuItemId,
			Quantity:   quantity,
		})

		fmt.Println("Menu item added to order successfully.")
	}

	orderID, err := h.apiClient.PostOrder(restaurantId, orderItems, token)
	if err != nil {
		fmt.Println("Error while creating order:", err)
		return
	}

	fmt.Printf("Order created successfully with ID: %d\n", orderID)
}

func (h *Handlers) HandleAddMenuItemToOrder(token string) {
	var orderId int
	var menuItemId int
	var quantity int

	fmt.Println("Enter Order ID:")
	fmt.Scanln(&orderId)
	fmt.Println("Enter Menu Item ID to add to order:")
	fmt.Scanln(&menuItemId)
	fmt.Println("Enter quantity:")
	fmt.Scanln(&quantity)

	err := h.apiClient.PostItemToOrder(orderId, menuItemId, quantity, token)
	if err != nil {
		fmt.Println("Error while adding item to order:", err)
		return
	}
	fmt.Println("Menu item added to order successfully.")
}

func (h *Handlers) HandlePlaceOrderAndGetBill(token string) {
	var orderId int

	fmt.Println("Enter Order ID to place:")
	fmt.Scanln(&orderId)

	invoiceId, err := h.apiClient.PostCreateInvoice(orderId, token)
	if err != nil {
		fmt.Println("Error while placing order:", err)
		return
	}

	bill, err := h.apiClient.GetInvoiceById(invoiceId, token)
	if err != nil {
		fmt.Println("Error while fetching invoice:", err)
		return
	}

	fmt.Printf("Order placed successfully.\nInvoice ID: %d, Amount: %.2f, Tax: %.2f, Total to Pay: %.2f, Payment Status: %s\n", bill.ID, bill.Total, bill.Tax, bill.ToPay, bill.PaymentStatus)
}

func (h *Handlers) HandlePayBill(token string) {
	var invoiceId int
	var amount float64

	fmt.Println("Enter Invoice ID to pay:")
	fmt.Scanln(&invoiceId)
	fmt.Println("Enter amount to pay:")
	fmt.Scanln(&amount)

	err := h.apiClient.PostPayInvoice(invoiceId, amount, token)
	if err != nil {
		fmt.Println("Error while paying invoice:", err)
		return
	}

	fmt.Println("Invoice paid successfully.")
}

func (h *Handlers) HandleGetInvoiceById(token string) {
	var invoiceId int

	fmt.Println("Enter Invoice ID to pay:")
	fmt.Scanln(&invoiceId)

	invoice, err := h.apiClient.GetInvoiceById(invoiceId, token)

	if err != nil {
		fmt.Println("Error while fetching invoice: ", err)
		return
	}

	fmt.Printf("Invoice fetched successfully.\nInvoice ID: %d, Amount: %.2f, Tax: %.2f, Total to Pay: %.2f, Payment Status: %s\n", invoice.ID, invoice.Total, invoice.Tax, invoice.ToPay, invoice.PaymentStatus)
}
