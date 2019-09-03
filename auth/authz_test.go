package auth

import (
	"testing"
)

type TestRole uint

const (
	REGULAR TestRole = iota
	MANAGER
	ADMIN
)

var allRoles = []TestRole{REGULAR, MANAGER, ADMIN}
var stringRoles = []string{"regular", "manager", "admin"}

func (r TestRole) String() string {
	return stringRoles[r]
}

func (t TestRole) Role() string {
	return t.String()
}

type TestResource uint

const (
	USER TestResource = iota
	LOSER
)

var allResources = []TestResource{USER, LOSER}
var stringResources = []string{"user", "loser"}

func (r TestResource) String() string {
	return stringRoles[r]
}

func (t TestResource) Resource() string {
	return t.String()
}

func TestOnlyReadForOneRole(t *testing.T) {
	// Arrange
	authorizer := NewAuthzService()
	authorizer.AddPermission(REGULAR, USER, Read)

	// Act & Assert
	assert(t, authorizer.Allowed(REGULAR, USER, Read))
	assert(t, !authorizer.Allowed(REGULAR, USER, Create))
	assert(t, !authorizer.Allowed(REGULAR, USER, Delete))
	assert(t, !authorizer.Allowed(REGULAR, USER, Update))
}

func TestCreateReadTwoRoles(t *testing.T) {
	// Arrange
	authorizer := NewAuthzService()
	authorizer.AddPermission(REGULAR, USER, Read)
	authorizer.AddPermission(MANAGER, USER, Create, Read, Update)

	// Act & Assert
	assert(t, authorizer.Allowed(REGULAR, USER, Read))
	assert(t, !authorizer.Allowed(REGULAR, USER, Create))
	assert(t, !authorizer.Allowed(REGULAR, USER, Delete))
	assert(t, !authorizer.Allowed(REGULAR, USER, Update))

	assert(t, authorizer.Allowed(MANAGER, USER, Create))
	assert(t, authorizer.Allowed(MANAGER, USER, Read))
	assert(t, authorizer.Allowed(MANAGER, USER, Update))
	assert(t, !authorizer.Allowed(MANAGER, USER, Delete))
}

func TestTwoResources(t *testing.T) {
	// Arrange
	authorizer := NewAuthzService()
	authorizer.AddPermission(REGULAR, USER, Read)
	authorizer.AddPermission(REGULAR, LOSER, Create)

	// Act & Assert
	assert(t, authorizer.Allowed(REGULAR, USER, Read))
	assert(t, !authorizer.Allowed(REGULAR, USER, Create))
	assert(t, !authorizer.Allowed(REGULAR, USER, Delete))
	assert(t, !authorizer.Allowed(REGULAR, USER, Update))
	assert(t, authorizer.Allowed(REGULAR, LOSER, Create))
	assert(t, !authorizer.Allowed(REGULAR, LOSER, Update))
	assert(t, !authorizer.Allowed(REGULAR, LOSER, Read))
	assert(t, !authorizer.Allowed(REGULAR, LOSER, Delete))
}

func assert(t *testing.T, expr bool) {
	t.Helper()

	if !expr {
		t.Error("Expression is false")
	}
}
