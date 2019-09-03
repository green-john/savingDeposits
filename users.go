package savingDeposits

type Role uint

const (
	REGULAR Role = iota
	MANAGER
	ADMIN
)

var AllRoles = []Role{REGULAR, MANAGER, ADMIN}
var stringRoles = []string{"regular", "manager", "admin"}

func (r Role) String() string {
	return stringRoles[r]
}

func (r Role) Role() string {
	return r.String()
}

type Resource uint

const (
	USERS Resource = iota
)

var AllResources = []Resource{USERS}
var stringResources = []string{"users"}

func (r Resource) String() string {
	return stringRoles[r]
}

func (r Resource) Resource() string {
	return r.String()
}

func ResourceFromString(s string) Resource {
	m := map[string]Resource{
		"users": USERS,
	}

	return m[s]
}

func RoleFromString(s string) Role {
	m := map[string]Role{
		"regular": REGULAR,
		"manager": MANAGER,
		"admin":   ADMIN,
	}

	return m[s]
}

type User struct {
	// Primary key
	ID uid `gorm:"primary_key" json:"id"`

	// Username
	Username string `gorm:"unique" json:"username"`

	// Password hash. Not included in json responses
	PasswordHash string `json:"-"`

	// Role
	Role string `json:"role"`
}

type UserSession struct {
	// Primary key
	ID uint `gorm:"primary_key"`

	// Generated token
	Token string

	// User associated to this session
	UserID uint
	User   User
}

type UserCreateInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserCreateOutput struct {
	User
}

type UserReadInput struct {
	// ID to lookup the apartment
	Id string
}

type UserReadOutput struct {
	User
}

type UserAllInput struct {
}

type UserAllOutput struct {
	Users []User
}

func (o *UserAllOutput) Public() interface{} {
	return o.Users
}

type UserUpdateInput struct {
	Id       string `json:"-"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserUpdateOutput struct {
	User
}

type UserDeleteInput struct {
	Id string
}

type UserDeleteOutput struct {
	Message string `json:"message"`
}

type UserService interface {
	Create(UserCreateInput) (*UserCreateOutput, error)
	Read(UserReadInput) (*UserReadOutput, error)
	All(UserAllInput) (*UserAllOutput, error)
	Update(UserUpdateInput) (*UserUpdateOutput, error)
	Delete(UserDeleteInput) (*UserDeleteOutput, error)
}
