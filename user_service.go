package main

import (
	"batch/panera"
	"context"
	"fmt"
)

const UserResolverID = "UserService"

var users = []string{
	"Chris",
	"Mike",
	"Lucy",
	"Katie",
	"Jamies",
}

type User struct {
	Name string
}

type UserService struct {
	panera.Resolver
}

func (u UserService) ID() string {
	return UserResolverID
}

func (u UserService) Resolve(ctx context.Context, queries map[int]any) map[int]any {
	fmt.Println("Detected", len(queries), "queries to the same service")
	results := make(map[int]any, len(queries))

	for id, query := range queries {
		requestedID := query.(int)
		result := User{
			Name: users[requestedID],
		}

		results[id] = result
	}

	return results
}

func (u UserService) Fetch(id int) panera.Node[User] {
	return panera.NewBatchQueryNode[int, User](
		func(_ context.Context) int {
			return id
		},
		UserResolverID,
	)
}
