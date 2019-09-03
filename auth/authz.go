package auth

type Permission uint

const (
	Create Permission = iota
	Read
	Update
	Delete
)

type Role interface {
	Role() string
}

type Resource interface {
	Resource() string
}

type AuthzService struct {
	// This is a map of roles to permissions per resource
	//
	//		role1 -> {
	//			resource1 -> [Read, Write]
	//			resource2 -> [Read, Write]
	//		},
	//		role2 -> {
	//			resource1 -> [Create, Delete]
	//			resource3 -> [Read, Update]
	//		}
	//
	perm map[string]map[string][]Permission
}

func (a *AuthzService) AddPermission(role Role, resource Resource, permissions ...Permission) {
	if _, ok := a.perm[role.Role()]; !ok {
		a.perm[role.Role()] = make(map[string][]Permission)
	}

	if _, ok := a.perm[role.Role()][resource.Resource()]; !ok {
		a.perm[role.Role()][resource.Resource()] = make([]Permission, 0)
	}

	for _, p := range permissions {
		if !containsPerm(a.perm[role.Role()][resource.Resource()], p) {
			a.perm[role.Role()][resource.Resource()] = append(a.perm[role.Role()][resource.Resource()], p)
		}
	}
}

func (a *AuthzService) Allowed(role Role, resource Resource, permission Permission) bool {
	// If role does not exist, return false
	if _, ok := a.perm[role.Role()]; !ok {
		return false
	}

	// Same for resource
	if _, ok := a.perm[role.Role()][resource.Resource()]; !ok {
		return false
	}

	return containsPerm(a.perm[role.Role()][resource.Resource()], permission)
}

func NewAuthzService() *AuthzService {
	p := make(map[string]map[string][]Permission)
	return &AuthzService{p}
}

func containsPerm(permissions []Permission, perm Permission) bool {
	for _, elt := range permissions {
		if elt == perm {
			return true
		}
	}

	return false
}
