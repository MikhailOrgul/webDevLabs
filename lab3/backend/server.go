package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
}

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	Role         string `json:"role"`
}

type CartItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type Cart struct {
	UserID int        `json:"userId"`
	Items  []CartItem `json:"items"`
}

type Store struct {
	Products []Product `json:"products"`
	Users    []User    `json:"users"`
	Carts    []Cart    `json:"carts"`
}

type ProductListResponse struct {
	Items      []Product `json:"items"`
	Page       int       `json:"page"`
	TotalPages int       `json:"totalPages"`
}

type CartItemView struct {
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	LineTotal float64 `json:"lineTotal"`
}

type CartResponse struct {
	Items []CartItemView `json:"items"`
	Total float64        `json:"total"`
}

const storeFilePath = "data/store.json"

var storeLock sync.Mutex

func main() {
	e := echo.New()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("very-secret-key"))))

	e.Static("/public", "public")
	e.Static("/assets", "public/assets")

	if err := ensureStoreFile(); err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/testget", func(c echo.Context) error {
		name := c.QueryParam("name")
		return c.String(http.StatusOK, "Привет, "+name)
	})

	e.POST("/testpost", func(c echo.Context) error {
		jsonMap := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
		if err != nil {
			return err
		}

		fmt.Println(jsonMap)
		name := jsonMap["name"].(string)
		v := map[string]interface{}{
			"response": "Привет, " + name,
		}

		return c.JSON(http.StatusOK, v)
	})

	e.POST("/api/register", func(c echo.Context) error {
		var payload struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректные данные."})
		}
		if payload.Name == "" || payload.Email == "" || payload.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Заполните все поля."})
		}

		store, err := loadStore()
		if err != nil {
			return err
		}
		if findUserByEmail(store.Users, payload.Email) != nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Пользователь уже существует."})
		}

		user := User{
			ID:           nextUserID(store.Users),
			Name:         payload.Name,
			Email:        payload.Email,
			PasswordHash: hashPassword(payload.Password),
			Role:         "user",
		}
		store.Users = append(store.Users, user)
		if err := saveStore(store); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, sanitizeUser(user))
	})

	e.POST("/api/login", func(c echo.Context) error {
		var payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректные данные."})
		}

		store, err := loadStore()
		if err != nil {
			return err
		}
		user := findUserByEmail(store.Users, payload.Email)
		if user == nil || user.PasswordHash != hashPassword(payload.Password) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверная почта или пароль."})
		}

		sess, _ := session.Get("session", c)
		sess.Values["userId"] = user.ID
		sess.Save(c.Request(), c.Response())

		return c.JSON(http.StatusOK, sanitizeUser(*user))
	})

	e.POST("/api/logout", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options.MaxAge = -1
		sess.Save(c.Request(), c.Response())
		return c.NoContent(http.StatusNoContent)
	})

	e.GET("/api/me", func(c echo.Context) error {
		user, err := getSessionUser(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Нет авторизации."})
		}
		return c.JSON(http.StatusOK, sanitizeUser(*user))
	})

	e.GET("/api/products", func(c echo.Context) error {
		store, err := loadStore()
		if err != nil {
			return err
		}

		page := parsePositiveInt(c.QueryParam("page"), 1)
		limit := parsePositiveInt(c.QueryParam("limit"), 3)
		totalPages := (len(store.Products) + limit - 1) / limit
		start := (page - 1) * limit
		if start > len(store.Products) {
			start = len(store.Products)
		}
		end := start + limit
		if end > len(store.Products) {
			end = len(store.Products)
		}

		response := ProductListResponse{
			Items:      store.Products[start:end],
			Page:       page,
			TotalPages: totalPages,
		}
		return c.JSON(http.StatusOK, response)
	})

	e.POST("/api/products", func(c echo.Context) error {
		if _, err := requireRole(c, "admin"); err != nil {
			return err
		}
		var payload Product
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректные данные."})
		}
		if payload.Name == "" || payload.Description == "" || payload.Price <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Заполните все поля."})
		}

		store, err := loadStore()
		if err != nil {
			return err
		}
		payload.ID = nextProductID(store.Products)
		store.Products = append(store.Products, payload)
		if err := saveStore(store); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, payload)
	})

	e.PUT("/api/products/:id", func(c echo.Context) error {
		if _, err := requireRole(c, "admin"); err != nil {
			return err
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный ID."})
		}
		var payload Product
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректные данные."})
		}
		if payload.Name == "" || payload.Description == "" || payload.Price <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Заполните все поля."})
		}

		store, err := loadStore()
		if err != nil {
			return err
		}
		updated := false
		for i := range store.Products {
			if store.Products[i].ID == id {
				payload.ID = id
				store.Products[i] = payload
				updated = true
				break
			}
		}
		if !updated {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Товар не найден."})
		}
		if err := saveStore(store); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, payload)
	})

	e.DELETE("/api/products/:id", func(c echo.Context) error {
		if _, err := requireRole(c, "admin"); err != nil {
			return err
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный ID."})
		}
		store, err := loadStore()
		if err != nil {
			return err
		}
		index := -1
		for i, product := range store.Products {
			if product.ID == id {
				index = i
				break
			}
		}
		if index == -1 {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Товар не найден."})
		}
		store.Products = append(store.Products[:index], store.Products[index+1:]...)
		if err := saveStore(store); err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	})

	e.GET("/api/cart", func(c echo.Context) error {
		user, err := requireAuth(c)
		if err != nil {
			return err
		}
		store, err := loadStore()
		if err != nil {
			return err
		}
		cart := findCart(store.Carts, user.ID)
		productByID := make(map[int]Product, len(store.Products))
		for _, product := range store.Products {
			productByID[product.ID] = product
		}

		response := CartResponse{}
		if cart != nil {
			for _, item := range cart.Items {
				product, ok := productByID[item.ProductID]
				if !ok {
					continue
				}
				lineTotal := product.Price * float64(item.Quantity)
				response.Items = append(response.Items, CartItemView{
					Product:   product,
					Quantity:  item.Quantity,
					LineTotal: lineTotal,
				})
				response.Total += lineTotal
			}
		}
		return c.JSON(http.StatusOK, response)
	})

	e.POST("/api/cart", func(c echo.Context) error {
		user, err := requireAuth(c)
		if err != nil {
			return err
		}
		var payload struct {
			ProductID int `json:"productId"`
			Quantity  int `json:"quantity"`
		}
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректные данные."})
		}
		if payload.ProductID <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный товар."})
		}
		if payload.Quantity <= 0 {
			payload.Quantity = 1
		}

		store, err := loadStore()
		if err != nil {
			return err
		}
		if findProductByID(store.Products, payload.ProductID) == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Товар не найден."})
		}

		cart := findCart(store.Carts, user.ID)
		if cart == nil {
			store.Carts = append(store.Carts, Cart{
				UserID: user.ID,
				Items: []CartItem{{ProductID: payload.ProductID, Quantity: payload.Quantity}},
			})
		} else {
			updated := false
			for i := range cart.Items {
				if cart.Items[i].ProductID == payload.ProductID {
					cart.Items[i].Quantity += payload.Quantity
					updated = true
					break
				}
			}
			if !updated {
				cart.Items = append(cart.Items, CartItem{ProductID: payload.ProductID, Quantity: payload.Quantity})
			}
		}

		if err := saveStore(store); err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	})

	e.DELETE("/api/cart/:id", func(c echo.Context) error {
		user, err := requireAuth(c)
		if err != nil {
			return err
		}
		productID, err := strconv.Atoi(c.Param("id"))
		if err != nil || productID <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный товар."})
		}
		store, err := loadStore()
		if err != nil {
			return err
		}
		cart := findCart(store.Carts, user.ID)
		if cart == nil {
			return c.NoContent(http.StatusNoContent)
		}
		newItems := cart.Items[:0]
		for _, item := range cart.Items {
			if item.ProductID != productID {
				newItems = append(newItems, item)
			}
		}
		cart.Items = newItems
		if err := saveStore(store); err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	})

	e.POST("/api/cart/clear", func(c echo.Context) error {
		user, err := requireAuth(c)
		if err != nil {
			return err
		}
		store, err := loadStore()
		if err != nil {
			return err
		}
		cart := findCart(store.Carts, user.ID)
		if cart != nil {
			cart.Items = nil
			if err := saveStore(store); err != nil {
				return err
			}
		}
		return c.NoContent(http.StatusNoContent)
	})

	e.GET("*", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func loadStore() (Store, error) {
	storeLock.Lock()
	defer storeLock.Unlock()

	var store Store
	data, err := os.ReadFile(storeFilePath)
	if err != nil {
		return store, err
	}
	if err := json.Unmarshal(data, &store); err != nil {
		return store, err
	}
	return store, nil
}

func saveStore(store Store) error {
	storeLock.Lock()
	defer storeLock.Unlock()

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(storeFilePath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(storeFilePath, data, 0o644)
}

func ensureStoreFile() error {
	if _, err := os.Stat(storeFilePath); err == nil {
		return nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	defaultStore := Store{
		Products: defaultProducts(),
		Users: []User{
			{
				ID:           1,
				Name:         "Администратор",
				Email:        "admin@shop.local",
				PasswordHash: hashPassword("admin123"),
				Role:         "admin",
			},
		},
	}
	return saveStore(defaultStore)
}

func defaultProducts() []Product {
	return []Product{
		{
			ID:          1,
			Name:        "Масляный фильтр ВАЗ",
			Description: "Качественный масляный фильтр для автомобилей ВАЗ. Обеспечивает надежную очистку масла и продлевает срок службы двигателя.",
			Price:       450,
			Image:       "/assets/img/фильтр.jpg",
		},
		{
			ID:          2,
			Name:        "Тормозные колодки ВАЗ",
			Description: "Комплект передних тормозных колодок для автомобилей ВАЗ. Обеспечивает эффективное торможение и устойчивость на дороге.",
			Price:       1200,
			Image:       "/assets/img/колодки.jpg",
		},
		{
			ID:          3,
			Name:        "Аккумулятор 60 А·ч",
			Description: "Надежный аккумулятор емкостью 60 А·ч для легковых автомобилей. Обеспечивает уверенный запуск двигателя в любое время года.",
			Price:       5800,
			Image:       "/assets/img/акум.jpg",
		},
		{
			ID:          4,
			Name:        "Масло моторное 10W-40",
			Description: "Полусинтетическое масло 10W-40 для стабильной работы двигателя в разных режимах.",
			Price:       2100,
			Image:       "/assets/img/t1.jpeg",
		},
		{
			ID:          5,
			Name:        "Свечи зажигания комплект",
			Description: "Комплект свечей зажигания с увеличенным ресурсом и стабильной искрой.",
			Price:       950,
			Image:       "/assets/img/t2.jpeg",
		},
		{
			ID:          6,
			Name:        "Ремень ГРМ",
			Description: "Надежный ремень ГРМ для автомобилей ВАЗ с повышенной износостойкостью.",
			Price:       1350,
			Image:       "/assets/img/t3.jpeg",
		},
		{
			ID:          7,
			Name:        "Фильтр салона",
			Description: "Салонный фильтр для очистки воздуха в автомобиле от пыли и аллергенов.",
			Price:       650,
			Image:       "/assets/img/p1.jpeg",
		},
	}
}

func sanitizeUser(user User) map[string]interface{} {
	return map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}
}

func nextProductID(products []Product) int {
	maxID := 0
	for _, product := range products {
		if product.ID > maxID {
			maxID = product.ID
		}
	}
	return maxID + 1
}

func nextUserID(users []User) int {
	maxID := 0
	for _, user := range users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}
	return maxID + 1
}

func findProductByID(products []Product, id int) *Product {
	for i := range products {
		if products[i].ID == id {
			return &products[i]
		}
	}
	return nil
}

func findUserByEmail(users []User, email string) *User {
	for i := range users {
		if users[i].Email == email {
			return &users[i]
		}
	}
	return nil
}

func findUserByID(users []User, id int) *User {
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

func findCart(carts []Cart, userID int) *Cart {
	for i := range carts {
		if carts[i].UserID == userID {
			return &carts[i]
		}
	}
	return nil
}

func getSessionUser(c echo.Context) (*User, error) {
	userID, ok := getSessionUserID(c)
	if !ok {
		return nil, errors.New("no session")
	}
	store, err := loadStore()
	if err != nil {
		return nil, err
	}
	user := findUserByID(store.Users, userID)
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func getSessionUserID(c echo.Context) (int, bool) {
	sess, _ := session.Get("session", c)
	rawID, ok := sess.Values["userId"]
	if !ok {
		return 0, false
	}
	switch value := rawID.(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}

func requireAuth(c echo.Context) (*User, error) {
	user, err := getSessionUser(c)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"error": "Требуется авторизация."})
	}
	return user, nil
}

func requireRole(c echo.Context, role string) (*User, error) {
	user, err := requireAuth(c)
	if err != nil {
		return nil, err
	}
	if user.Role != role {
		return nil, c.JSON(http.StatusForbidden, map[string]string{"error": "Недостаточно прав."})
	}
	return user, nil
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
