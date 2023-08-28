package main

import (
	"context"
	"fmt"
	"strings"

	"batch/panera"
)

func main() {
	userService := UserService{}

	// This is an example of a non-trivial execution graph. FetchUser is intended to simulate
	// an RPC.
	//
	// FetchUser(0) --> FetchUser(1)  |
	// FetchUser(1) --> FetchUser(2)  |--> Merge results into a string
	// FetchUser(2) --------------->  |
	//
	// In this scenario we see two groups of RPCs to the same service. We are able to identify
	// both groups of calls and batch them appropriately.

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

	result, _ := panera.ExecuteGraph[string](
		context.Background(),
		users,
		map[string]panera.Resolver{
			userService.ID(): userService,
		},
	)

	fmt.Println(result)
}
