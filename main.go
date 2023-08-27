package main

import (
	"context"
	"fmt"
	"strings"
)

func main() {
	userService := UserService{}
	userResolver := UserResolver{}

	users := NewTransformNode[[]User, string](
		NewListNode([]Node[User]{
			NewFlatMapNode(
				userService.Fetch(0),
				func(result User) Node[User] {
					return userService.Fetch(1)
				},
			),
			NewFlatMapNode(
				userService.Fetch(1),
				func(result User) Node[User] {
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

	result := ExecuteGraph[string](
		context.Background(),
		users,
		map[string]Resolver{
			userResolver.ID(): userResolver,
		},
	)

	fmt.Println(result)
}
