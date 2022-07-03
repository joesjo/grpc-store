package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/joesjo/grpc-store/shopinterface/graph/generated"
	"github.com/joesjo/grpc-store/shopinterface/graph/model"
	"github.com/joesjo/grpc-store/shopinterface/serviceclient"
)

func (r *mutationResolver) CreateItem(ctx context.Context, name string, quantity int) (*model.Item, error) {
	itemId, err := serviceclient.CreateItem(name, int32(quantity))
	if err != nil {
		return nil, err
	}
	item, err := serviceclient.GetItem(itemId)
	if err != nil {
		return nil, err
	}
	return &model.Item{
		ID:       item.GetId(),
		Name:     item.GetName(),
		Quantity: int(item.GetQuantity()),
	}, nil
}

func (r *mutationResolver) UpdateItem(ctx context.Context, id string, name *string, quantity *int) (*model.Item, error) {
	if name == nil || quantity == nil {
		return nil, fmt.Errorf("optional parameters are not supported yet")
	}
	err := serviceclient.UpdateItem(id, *name, int32(*quantity))
	if err != nil {
		return nil, err
	}
	item, err := serviceclient.GetItem(id)
	if err != nil {
		return nil, err
	}
	return &model.Item{
		ID:       item.GetId(),
		Name:     item.GetName(),
		Quantity: int(item.GetQuantity()),
	}, nil
}

func (r *mutationResolver) DeleteItem(ctx context.Context, id string) (bool, error) {
	err := serviceclient.DeleteItem(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) IncrementItem(ctx context.Context, input model.IncrementItem) (*model.Item, error) {
	err := serviceclient.StockItem(input.ID, int32(input.Quantity))
	if err != nil {
		return nil, err
	}
	item, err := serviceclient.GetItem(input.ID)
	if err != nil {
		return nil, err
	}
	return &model.Item{
		ID:       item.GetId(),
		Name:     item.GetName(),
		Quantity: int(item.GetQuantity()),
	}, nil
}

func (r *queryResolver) Items(ctx context.Context) ([]*model.Item, error) {
	itemArray, err := serviceclient.GetInventory()
	if err != nil {
		return nil, err
	}
	items := make([]*model.Item, len(itemArray))
	for i, item := range itemArray {
		items[i] = &model.Item{
			ID:       item.GetId(),
			Name:     item.GetName(),
			Quantity: int(item.GetQuantity()),
		}
	}
	return items, nil
}

func (r *queryResolver) Item(ctx context.Context, id string) (*model.Item, error) {
	item, err := serviceclient.GetItem(id)
	if err != nil {
		return nil, err
	}
	return &model.Item{
		ID:       item.GetId(),
		Name:     item.GetName(),
		Quantity: int(item.GetQuantity()),
	}, nil
}

func (r *queryResolver) FindItems(ctx context.Context, name string) ([]*model.Item, error) {
	itemArray, err := serviceclient.FindItems(name)
	if err != nil {
		return nil, err
	}
	items := make([]*model.Item, len(itemArray))
	for i, item := range itemArray {
		items[i] = &model.Item{
			ID:       item.GetId(),
			Name:     item.GetName(),
			Quantity: int(item.GetQuantity()),
		}
	}
	return items, nil
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
func (r *queryResolver) FindItem(ctx context.Context, name string) (*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}
