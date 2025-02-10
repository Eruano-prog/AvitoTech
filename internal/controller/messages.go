package controller

// AuthRequest defines model for AuthRequest.
type AuthRequest struct {
	// Password Пароль для аутентификации.
	Password string `json:"password"`

	// Username Имя пользователя для аутентификации.
	Username string `json:"username"`
}

// AuthResponse defines model for AuthResponse.
type AuthResponse struct {
	// Token JWT-токен для доступа к защищенным ресурсам.
	Token *string `json:"token,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Errors Сообщение об ошибке, описывающее проблему.
	Errors *string `json:"errors,omitempty"`
}

// InfoResponse defines model for InfoResponse.
type InfoResponse struct {
	CoinHistory *History `json:"coinHistory,omitempty"`

	// Coins Количество доступных монет.
	Coins     *int               `json:"coins,omitempty"`
	Inventory *[]InventoryRecord `json:"inventory,omitempty"`
}

type History struct {
	Received *[]ReceiveRecord `json:"received,omitempty"`
	Sent     *[]SendRecord    `json:"sent,omitempty"`
}

type InventoryRecord struct {
	// Quantity Количество предметов.
	Quantity *int `json:"quantity,omitempty"`

	// Type Тип предмета.
	Type *string `json:"type,omitempty"`
}

type ReceiveRecord struct {
	// Amount Количество полученных монет.
	Amount *int `json:"amount,omitempty"`

	// FromUser Имя пользователя, который отправил монеты.
	FromUser *string `json:"fromUser,omitempty"`
}

type SendRecord struct {
	// Amount Количество отправленных монет.
	Amount *int `json:"amount,omitempty"`

	// ToUser Имя пользователя, которому отправлены монеты.
	ToUser *string `json:"toUser,omitempty"`
}

// SendCoinRequest defines model for SendCoinRequest.
type SendCoinRequest struct {
	// Amount Количество монет, которые необходимо отправить.
	Amount int `json:"amount"`

	// ToUser Имя пользователя, которому нужно отправить монеты.
	ToUser string `json:"toUser"`
}

// PostApiAuthJSONRequestBody defines body for ApiAuth for application/json ContentType.
type PostApiAuthJSONRequestBody = AuthRequest

// PostApiSendCoinJSONRequestBody defines body for ApiSendCoin for application/json ContentType.
type PostApiSendCoinJSONRequestBody = SendCoinRequest
