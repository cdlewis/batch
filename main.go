package main

import (
	"fmt"
	"strings"
)

func main() {
	userService := UserService{}
	userResolver := UserResolver{}

	users := NewTransformNode[[]User, string](
		NewListNode([]Node[User]{
			userService.Fetch(0),
			userService.Fetch(1),
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
		users,
		map[string]Resolver{
			userResolver.ID(): userResolver,
		},
	)

	fmt.Println(result)
}
