package repository

import (
	"kuncie/graph/model"
	db "kuncie/internal/pkg/db/mysql"
	"log"

	"github.com/davecgh/go-spew/spew"
)

//CreateUser create's user
func CreateUser(user model.User) (int, error) {

	stmt, err := db.Db.Prepare("INSERT INTO users(name, email, phone_number, address) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	res, err := stmt.Exec(user.Name, user.Email, user.PhoneNumber, user.Address)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	defer stmt.Close()
	log.Println("Row inserted!!")
	return int(id), nil
}

//GetUserByID return user with respective id
func GetUserByID(id *int) (*model.User, error) {

	stmt, err := db.Db.Prepare("select * from users where id=?")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()
	var user model.User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Address)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &user, nil

}

//CreateProduct create's product
func CreateProduct(product model.Product) (int, error) {

	stmt, err := db.Db.Prepare("INSERT INTO products(sku, name, price, qty) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	res, err := stmt.Exec(product.Sku, product.Name, product.Price, product.Qty)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	defer stmt.Close()
	log.Println("Row inserted!!")
	return int(id), nil
}

//GetProductByID return product with respective id
func GetProductByID(id *int) (*model.Product, error) {

	stmt, err := db.Db.Prepare("select * from products where id=?")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()
	var product model.Product
	for rows.Next() {
		err = rows.Scan(&product.ID, &product.Sku, &product.Name, &product.Price, &product.Qty, &product.PromotionID)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &product, nil

}

//GetAllUsers returns all users
func GetAllUsers() ([]*model.User, error) {
	stmt, err := db.Db.Prepare("select * from users")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var users []*model.User
	for rows.Next() {
		var user model.User
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Address)
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer stmt.Close()
	defer rows.Close()

	return users, err
}

//GetAllProducts returns all products
func GetAllProducts() ([]*model.Product, error) {
	stmt, err := db.Db.Prepare("select * from products")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		rows.Scan(&product.ID, &product.Sku, &product.Name, &product.Price, &product.Qty, &product.PromotionID)
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer stmt.Close()
	defer rows.Close()

	return products, err
}

//GetTransactionByID return product with respective id
func GetTransactionByID(id *int) (*model.Transaction, error) {

	stmt, err := db.Db.Prepare("select t.id, t.grand_total, u.id from transactions t inner join users u where t.user_id = u.id and t.id = ? ;")
	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(id)
	var trans model.Transaction
	if row.Next() {
		err := row.Scan(&trans.ID, &trans.GrandTotal, &trans.UserID)
		if err != nil {
			return nil, err
		}
	}

	stmt, err = db.Db.Prepare("SELECT td.product_id, td.qty, td.sub_total FROM transaction_details td LEFT JOIN products p ON td.product_id = p.id WHERE td.transaction_id = ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	var details []*model.TransactionDetail
	for rows.Next() {
		var detail model.TransactionDetail
		rows.Scan(&detail.ProductID, &detail.Qty, &detail.SubTotal)
		details = append(details, &detail)

	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	transaction := &model.Transaction{
		ID:         trans.ID,
		GrandTotal: trans.GrandTotal,
		UserID:     trans.UserID,
		Details:    details,
	}
	defer row.Close()
	defer stmt.Close()
	return transaction, nil

}

//CreateTransaction create's transaction
func CreateTransaction(transaction model.Transaction) (int, error) {

	stmt, err := db.Db.Prepare("INSERT INTO transactions(user_id, grand_total) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	res, err := stmt.Exec(transaction.UserID, 0)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	var grandTotal = 0
	var totalDiscount = 0
	for _, detail := range transaction.Details {
		product, err := GetProductByID(&detail.ProductID)

		if err != nil {
			return 0, err
		}

		if product.Qty < detail.Qty {
			return 0, nil
		}

		// var divQty = 0
		// var promoPrice = 0
		// var modQty = 0
		// var regularPrice = 0
		// var subTotalPrice = 0
		var discount = 0

		if product.PromotionID == 1 {

		} else if product.PromotionID == 2 {
			stmt, err := db.Db.Prepare("select * from promo_payless_rules ppr where ppr.promotion_id = ? ;")
			if err != nil {
				return 0, err
			}

			row, err := stmt.Query(product.PromotionID)
			var ppr model.PromoPaylessRule
			if row.Next() {
				err := row.Scan(&ppr.PromotionID, &ppr.RequirementQty, &ppr.PromoQty)
				if err != nil {
					return 0, err
				}
			}

			if detail.Qty >= ppr.RequirementQty {
				var divQty = detail.Qty / ppr.RequirementQty
				var promoPrice = divQty * ppr.PromoQty * product.Price
				var modQty = detail.Qty % ppr.RequirementQty
				var regularPrice = modQty * product.Price
				var subTotalPrice = detail.Qty * product.Price
				discount = subTotalPrice - (promoPrice + regularPrice)

				spew.Dump(divQty, promoPrice, modQty, regularPrice)
			}
		} else if product.PromotionID == 3 {
			stmt, err := db.Db.Prepare("select * from promo_discount_rules pdr where pdr.promotion_id = ? ;")
			if err != nil {
				return 0, err
			}

			row, err := stmt.Query(product.PromotionID)
			var pdr model.PromoDiscountRule
			if row.Next() {
				err := row.Scan(&pdr.PromotionID, &pdr.RequirementQty, &pdr.PercentageDiscount)
				if err != nil {
					return 0, err
				}
			}

			if detail.Qty >= pdr.RequirementQty {
				discount = (detail.Qty * product.Price) * pdr.PercentageDiscount / 100

				spew.Dump(detail.Qty, product.Price, pdr.PercentageDiscount)
			}
		}

		stmt, err = db.Db.Prepare("INSERT INTO transaction_details(transaction_id, product_id, price, qty, sub_total, discount) VALUES(?,?,?,?,?,?)")
		if err != nil {
			log.Fatal(err)
			return 0, err
		}

		var subTotal = detail.Qty * product.Price
		res, err = stmt.Exec(id, detail.ProductID, product.Price, detail.Qty, subTotal, discount)
		if err != nil {
			log.Fatal(err)
			return 0, err
		}

		var qty = product.Qty - detail.Qty
		stmt, err = db.Db.Prepare("UPDATE products SET qty=? WHERE id=?")
		if err != nil {
			log.Fatal(err)
			return 0, err
		}

		res, err = stmt.Exec(qty, product.ID)
		if err != nil {
			log.Fatal(err)
			return 0, err
		}

		grandTotal = grandTotal + subTotal
		totalDiscount = totalDiscount + discount
	}

	grandTotal = grandTotal - totalDiscount

	stmt, err = db.Db.Prepare("UPDATE transactions SET grand_total=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	res, err = stmt.Exec(grandTotal, id)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	defer stmt.Close()
	log.Println("Row inserted!!")
	return int(id), nil
}

func CheckPromoPaylessRule(product model.Product) (discount float32) {
	stmt, err := db.Db.Prepare("select * from promo_payless_rules ppr where ppr.promotion_id = ? ;")
	if err != nil {
		return 0
	}

	row, err := stmt.Query(2)
	var ppr model.PromoPaylessRule
	if row.Next() {
		err := row.Scan(&ppr.PromotionID, &ppr.RequirementQty, &ppr.PromoQty)
		if err != nil {
			return 0
		}
	}

	if product.Qty >= ppr.RequirementQty {
		var divQty = product.Qty / ppr.RequirementQty
		var promoPrice = divQty * ppr.PromoQty * product.Price
		var modQty = product.Qty % ppr.RequirementQty
		var regularPrice = modQty * product.Price

		spew.Dump(promoPrice, regularPrice)
	}

	return 0
}
