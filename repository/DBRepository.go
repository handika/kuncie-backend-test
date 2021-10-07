package repository

import (
	"kuncie/graph/model"
	db "kuncie/internal/pkg/db/mysql"
	"log"
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
	var freeItems map[int]int
	freeItems = map[int]int{}
	var mainProductId = 0
	var freeProductId = 0
	for _, detail := range transaction.Details {
		product, err := GetProductByID(&detail.ProductID)

		if err != nil {
			return 0, err
		}

		if product.Qty < detail.Qty {
			return 0, nil
		}

		var discount = 0
		freeItems[product.ID] = detail.Qty

		if product.PromotionID == 1 {
			stmt, err := db.Db.Prepare("select * from promo_free_item_rules pfir where pfir.promotion_id = ? ;")
			if err != nil {
				return 0, err
			}

			row, err := stmt.Query(product.PromotionID)
			var pfir model.PromoFreeItemRule
			if row.Next() {
				err := row.Scan(&pfir.PromotionID, &pfir.FreeProductID)
				if err != nil {
					return 0, err
				}
			}

			mainProductId = product.ID
			freeProductId = pfir.FreeProductID
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

	var limitBuy, limitGet int
	for key, item := range freeItems {
		if key == mainProductId {
			limitBuy = item
		}

		if key == freeProductId {
			limitGet = item
		}
	}

	if limitBuy >= limitGet {
		stmt, err := db.Db.Prepare("select id, price from products where id = ? ;")
		if err != nil {
			return 0, err
		}

		row, err := stmt.Query(freeProductId)
		var prod model.Product
		if row.Next() {
			err := row.Scan(&prod.ID, &prod.Price)
			if err != nil {
				return 0, err
			}
		}

		stmt, err = db.Db.Prepare("UPDATE transaction_details SET discount=? WHERE transaction_id=? and product_id=?")
		if err != nil {
			log.Fatal(err)
			return 0, err
		}

		var discount = limitGet * prod.Price
		res, err = stmt.Exec(discount, id, freeProductId)
		if err != nil {
			log.Fatal(err)
			return 0, err
		}

		totalDiscount = totalDiscount + discount
	} else {
		stmt, err := db.Db.Prepare("select id, price from products where id = ? ;")
		if err != nil {
			return 0, err
		}

		row, err := stmt.Query(freeProductId)
		var prod model.Product
		if row.Next() {
			err := row.Scan(&prod.ID, &prod.Price)
			if err != nil {
				return 0, err
			}
		}

		stmt, err = db.Db.Prepare("UPDATE transaction_details SET discount=? WHERE transaction_id=? and product_id=?")
		if err != nil {
			log.Fatal(err)
			return 0, err
		}

		var discount = limitBuy * prod.Price
		res, err = stmt.Exec(limitBuy*prod.Price, id, freeProductId)
		if err != nil {
			log.Fatal(err)
			return 0, err
		}

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
