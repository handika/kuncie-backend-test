package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"kuncie/graph/generated"
	"kuncie/graph/model"
	"kuncie/repository"
)

func (r *mutationResolver) CreateUser(ctx context.Context, name string, email string, phoneNumber string, address string) (*model.User, error) {
	var user model.User
	user.Name = name
	user.Email = email
	user.PhoneNumber = phoneNumber
	user.Address = address
	id, err := repository.CreateUser(user)
	if err != nil {
		return nil, err
	} else {
		return &model.User{ID: id, Name: user.Name, Email: user.Email, PhoneNumber: user.PhoneNumber, Address: user.Address}, nil
	}
}

func (r *mutationResolver) CreateProduct(ctx context.Context, sku string, name string, price int, qty int) (*model.Product, error) {
	var product model.Product
	product.Sku = sku
	product.Name = name
	product.Price = price
	product.Qty = qty
	id, err := repository.CreateProduct(product)
	if err != nil {
		return nil, err
	}
	createdProduct, _ := repository.GetProductByID(&id)
	return createdProduct, nil
}

func (r *mutationResolver) CreateTransaction(ctx context.Context, input model.TransactionInput) (*model.Transaction, error) {
	transaction := model.Transaction{
		UserID:  input.UserID,
		Details: mapItemsFromInput(input.Details),
	}

	id, err := repository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	createdTransaction, _ := repository.GetTransactionByID(&id)

	return createdTransaction, nil
}

func (r *queryResolver) UserByID(ctx context.Context, id *int) (*model.User, error) {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (r *queryResolver) AllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := repository.GetAllUsers()
	if err != nil {
		return nil, err
	} else {
		return users, err
	}
}

func (r *queryResolver) ProductByID(ctx context.Context, id *int) (*model.Product, error) {
	product, err := repository.GetProductByID(id)
	if err != nil {
		return nil, err
	} else {
		return product, nil
	}
}

func (r *queryResolver) AllProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := repository.GetAllProducts()
	if err != nil {
		return nil, err
	} else {
		return products, err
	}
}

func (r *queryResolver) TransactionByID(ctx context.Context, id *int) (*model.Transaction, error) {
	transaction, err := repository.GetTransactionByID(id)
	if err != nil {
		return nil, err
	} else {
		return transaction, nil
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func mapItemsFromInput(transactionDetailInput []*model.TransactionDetailInput) []*model.TransactionDetail {
	var items []*model.TransactionDetail
	for _, itemInput := range transactionDetailInput {
		items = append(items, &model.TransactionDetail{
			ProductID: itemInput.ProductID,
			Qty:       itemInput.Qty,
		})
	}
	return items
}
