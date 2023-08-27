package main

import (
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

type UserResolver struct {
	Resolver
}

func (u UserResolver) ID() string {
	return UserResolverID
}

func (u UserResolver) Resolve(nodeIDs []int, taskManager *TaskManager) {
	fmt.Println("Detected", len(nodeIDs), "queries to the same service")
	for _, id := range nodeIDs {
		node := taskManager.GetTask(id).(*customNode)

		requestedID := node.userID
		result := User{
			Name: users[requestedID],
		}

		node.InjectResult(result)
		fmt.Println("FINISHED", id)
		taskManager.FinishTask(id)
	}
}

type UserService struct{}
type User struct {
	Name string
}

func (u UserService) Fetch(id int) Node[User] {
	return newCustomNode(id)
}

type customNode struct {
	BatchableNode

	userID     int
	isResolved bool
	value      User
}

func newCustomNode(userID int) Node[User] {
	return &customNode{userID: userID}
}

func (v *customNode) GetValue() User {
	return v.value
}

func (v *customNode) IsResolved() bool {
	return v.isResolved
}

func (v *customNode) GetAnyResolvables() []AnyNode {
	return []AnyNode{}
}

func (v *customNode) Run(_ context.Context) any {
	panic("we should batch this -- you screwed up")
}

func (v *customNode) ResolverID() string {
	return UserResolverID
}

func (v *customNode) InjectResult(user User) {
	v.value = user
	v.isResolved = true
}
