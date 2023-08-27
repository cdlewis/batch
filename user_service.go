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

type UserResolver struct {
	panera.Resolver
}

func (u UserResolver) ID() string {
	return UserResolverID
}

func (u UserResolver) Resolve(ctx context.Context, nodeIDs []int, taskManager panera.TaskManager) {
	fmt.Println("Detected", len(nodeIDs), "queries to the same service")
	for _, id := range nodeIDs {
		node := taskManager.GetTask(id).(*customNode)

		requestedID := node.userID
		result := User{
			Name: users[requestedID],
		}

		node.InjectResult(result)
		taskManager.FinishTask(id)
	}
}

type UserService struct{}
type User struct {
	Name string
}

func (u UserService) Fetch(id int) panera.Node[User] {
	return newCustomNode(id)
}

type customNode struct {
	panera.BatchableNode

	userID     int
	isResolved bool
	value      User
}

func newCustomNode(userID int) panera.Node[User] {
	return &customNode{userID: userID}
}

func (v *customNode) GetValue(ctx context.Context, id int) User {
	return v.value
}

func (v *customNode) IsResolved(ctx context.Context, id int) bool {
	return v.isResolved
}

func (v *customNode) GetChildren() []panera.AnyNode {
	return []panera.AnyNode{}
}

func (v *customNode) Run(_ context.Context, id int) any {
	panic("we should batch this -- you screwed up")
}

func (v *customNode) ResolverID() string {
	return UserResolverID
}

func (v *customNode) InjectResult(user User) {
	v.value = user
	v.isResolved = true
}
