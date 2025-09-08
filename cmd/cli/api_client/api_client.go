package apiclient

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type APIClient struct {
	client  *http.Client
	baseUrl string
}

func NewAPIClient(baseUrl string) *APIClient {
	return &APIClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseUrl: baseUrl,
	}
}

func (c *APIClient) GetRestaurants() ([]domain.Restaurant, error) {
	resp, err := c.client.Get(c.baseUrl + "/api/restaurants")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return nil, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return nil, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.GetRestaurantsResponse](resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error decoding response %w", err)
	}

	restaurants := []domain.Restaurant{}
	for _, restaurant := range response.Restaurants {
		restaurants = append(restaurants, dtos.NewRestaurant(restaurant))
	}
	return restaurants, nil
}

func (c *APIClient) PostRestaurants(name string, token string) (int, error) {
	buf := bytes.NewBuffer(nil)
	createReqDto := dtos.CreateRestaurantRequest{Name: name}
	if err := encodeJson(buf, createReqDto); err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/api/restaurants", buf)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusCreated {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return 0, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return 0, errors.New(errResp.Message)
	}

	defer resp.Body.Close()

	response, err := decodeResponse[dtos.CreateRestaurantResponse](resp.Body)
	if err != nil {
		return 0, err
	}

	return response.ID, nil
}

func (c *APIClient) GetUserById(id int) (*domain.User, error) {
	idStr := strconv.Itoa(id)
	req, err := http.NewRequest("GET", c.baseUrl+"/api/users/"+idStr, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, errors.New("unknown error occurred while doing request: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return nil, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return nil, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.GetUserResponse](resp.Body)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		ID:    response.UserID,
		Name:  response.Name,
		Email: response.Email,
		Role:  domain.UserRole(response.Role),
	}
	return &user, nil
}

func (c *APIClient) PostUser(name, email, password, role string) (int, error) {
	buf := bytes.NewBuffer(nil)
	createReqDto := dtos.CreateUserRequest{Name: name, Email: email, Password: password, Role: role}
	if err := encodeJson(buf, createReqDto); err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/api/users", buf)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return 0, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return 0, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.CreateUserResponse](resp.Body)
	if err != nil {
		return 0, err
	}

	return response.UserID, nil
}

func (c *APIClient) PostLogin(email, password string) (string, error) {
	buf := bytes.NewBuffer(nil)
	loginReqDto := dtos.LoginRequest{Email: email, Password: password}
	if err := encodeJson(buf, loginReqDto); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/api/auth/login", buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return "", errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return "", errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.LoginResponse](resp.Body)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (c *APIClient) PostLogout(token string) error {
	req, err := http.NewRequest("POST", c.baseUrl+"/api/auth/logout", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return errors.New(errResp.Message)
	}

	return nil
}

func (c *APIClient) GetMenuItems(restaurantId int) ([]domain.MenuItem, error) {
	restaurantIdStr := strconv.Itoa(restaurantId)
	req, err := http.NewRequest("GET", c.baseUrl+"/api/restaurants/"+restaurantIdStr+"/items", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return nil, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return nil, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.GetMenuItemsResponse](resp.Body)
	if err != nil {
		return nil, err
	}

	menuItems := []domain.MenuItem{}
	for _, item := range response.Items {
		menuItems = append(menuItems, domain.NewMenuItem(
			item.ID,
			item.Name,
			item.Price,
			item.Available,
			restaurantId,
		))
	}
	return menuItems, nil
}

func (c *APIClient) PostMenuItem(restaurantId int, name string, price float64, available bool, token string) (int, error) {
	buf := bytes.NewBuffer(nil)
	createReqDto := dtos.AddMenuItemRequest{Name: name, Price: price, Available: available}
	if err := encodeJson(buf, createReqDto); err != nil {
		return 0, err
	}

	restaurantIdStr := strconv.Itoa(restaurantId)
	req, err := http.NewRequest("POST", c.baseUrl+"/api/restaurants/"+restaurantIdStr+"/items", buf)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return 0, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return 0, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.AddMenuItemResponse](resp.Body)
	if err != nil {
		return 0, err
	}

	return response.ID, nil
}

func (c *APIClient) PatchMenuItemAvailability(menuItemId int, available bool, token string) error {
	buf := bytes.NewBuffer(nil)
	updateReqDto := dtos.UpdateMenuItemAvailabilityRequest{Available: available}
	if err := encodeJson(buf, updateReqDto); err != nil {
		return err
	}

	menuItemIdStr := strconv.Itoa(menuItemId)
	req, err := http.NewRequest("PATCH", c.baseUrl+"/api/items/"+menuItemIdStr, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return errors.New(errResp.Message)
	}

	return nil
}

func (c *APIClient) GetOrderById(orderId int, token string) (*domain.Order, error) {
	orderIdStr := strconv.Itoa(orderId)
	req, err := http.NewRequest("GET", c.baseUrl+"/api/orders/"+orderIdStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return nil, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return nil, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.GetOrderByIdResponse](resp.Body)
	if err != nil {
		return nil, err
	}

	orderItems := []domain.OrderItem{}
	for _, item := range response.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
		})
	}

	order := domain.Order{
		ID:           response.ID,
		CustomerID:   response.CustomerID,
		RestaurantID: response.RestaurantID,
		OrderItems:   orderItems,
	}

	return &order, nil
}

func (c *APIClient) PostOrder(restaurantId int, orderItems []domain.OrderItem, token string) (int, error) {
	buf := bytes.NewBuffer(nil)
	createReqDto := dtos.CreateOrderRequest{
		RestaurantID: restaurantId,
		OrderItems:   []dtos.OrderItemsDTO{},
	}
	for _, item := range orderItems {
		createReqDto.OrderItems = append(createReqDto.OrderItems, dtos.OrderItemsDTO{
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
		})
	}

	if err := encodeJson(buf, createReqDto); err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/api/orders", buf)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return 0, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return 0, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.CreateOrderResponse](resp.Body)
	if err != nil {
		return 0, err
	}

	return response.ID, nil
}

func (c *APIClient) PostItemToOrder(orderId, menuItemId, quantity int, token string) error {
	buf := bytes.NewBuffer(nil)
	addReqDto := dtos.AddOrderItemRequest{
		MenuItemID: menuItemId,
		Quantity:   quantity,
	}
	if err := encodeJson(buf, addReqDto); err != nil {
		return err
	}

	orderIdStr := strconv.Itoa(orderId)
	req, err := http.NewRequest("POST", c.baseUrl+"/api/orders/"+orderIdStr+"/items", buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return errors.New(errResp.Message)
	}

	return nil
}

func (c *APIClient) PostCreateInvoice(orderId int, token string) (int, error) {
	orderIdStr := strconv.Itoa(orderId)
	req, err := http.NewRequest("POST", c.baseUrl+"/api/orders/"+orderIdStr+"/invoices", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return 0, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return 0, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.InvoiceResponse](resp.Body)
	if err != nil {
		return 0, err
	}

	return response.ID, nil
}

func (c *APIClient) GetInvoiceById(invoiceId int, token string) (*dtos.InvoiceResponse, error) {
	invoiceIdStr := strconv.Itoa(invoiceId)
	req, err := http.NewRequest("GET", c.baseUrl+"/api/invoices/"+invoiceIdStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return nil, errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return nil, errors.New(errResp.Message)
	}

	response, err := decodeResponse[dtos.InvoiceResponse](resp.Body)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *APIClient) PostPayInvoice(invoiceId int, amount float64, token string) error {
	invoiceIdStr := strconv.Itoa(invoiceId)

	buf := bytes.NewBuffer(nil)
	payReqDto := dtos.PaymentRequest{Amount: amount}
	if err := encodeJson(buf, payReqDto); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/api/invoices/"+invoiceIdStr+"/pay", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp, err := decodeError(resp.Body)
		if err != nil {
			return errors.New("unknown error occurred while doing request: " + err.Error())
		}
		return errors.New(errResp.Message)
	}

	return nil
}
