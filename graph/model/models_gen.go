// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Product struct {
	ID          int    `json:"id"`
	Sku         string `json:"sku"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Qty         int    `json:"qty"`
	PromotionID int    `json:"promotionId"`
}

type PromoDiscountRule struct {
	PromotionID        int `json:"promotionId"`
	RequirementMinQty  int `json:"requirementMinQty"`
	PercentageDiscount int `json:"percentageDiscount"`
}

type PromoFreeItemRule struct {
	PromotionID   int `json:"promotionId"`
	FreeProductID int `json:"freeProductId"`
}

type PromoPaylessRule struct {
	PromotionID    int `json:"promotionId"`
	RequirementQty int `json:"requirementQty"`
	PromoQty       int `json:"promoQty"`
}

type Promotion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Transaction struct {
	ID         int                  `json:"id"`
	UserID     int                  `json:"userId"`
	GrandTotal int                  `json:"grandTotal"`
	Details    []*TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID        int `json:"id"`
	ProductID int `json:"productId"`
	Price     int `json:"price"`
	Qty       int `json:"qty"`
	SubTotal  int `json:"subTotal"`
	Discount  int `json:"discount"`
}

type TransactionDetailInput struct {
	ProductID int `json:"productId"`
	Qty       int `json:"qty"`
}

type TransactionInput struct {
	UserID     int                       `json:"userId"`
	GrandTotal int                       `json:"grandTotal"`
	Details    []*TransactionDetailInput `json:"details"`
}

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}
