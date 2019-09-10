package savingDeposits

import "errors"

var NotFoundError = errors.New("entity not found")
var NotAuthorizedError = errors.New("user not authorized to perform that action")
