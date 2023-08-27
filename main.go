package main

import (
	"context"
	"fmt"
	"strings"

	"batch/panera"
)

func main() {
	userService := UserService{}
	userResolver := UserResolver{}

	users := panera.NewTransformNode[[]User, string](
		panera.NewListNode([]panera.Node[User]{
			panera.NewFlatMapNode(
				userService.Fetch(0),
				func(result User) panera.Node[User] {
					return userService.Fetch(1)
				},
			),
			panera.NewFlatMapNode(
				userService.Fetch(1),
				func(result User) panera.Node[User] {
					return userService.Fetch(2)
				},
			),
			userService.Fetch(2),
		}),
		func(results []User) string {
			userNames := []string{}

			for _, i := range results {
				userNames = append(userNames, i.Name)
			}

			return strings.Join(userNames, ",")
		},
	)

	result := panera.ExecuteGraph[string](
		context.Background(),
		users,
		map[string]panera.Resolver{
			userResolver.ID(): userResolver,
		},
	)

	fmt.Println(result)
}
