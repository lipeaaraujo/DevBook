package users

import "errors"

var ErrNameEmpty = errors.New("User name can't be null or empty")
var ErrNicknameEmpty = errors.New("User nickname can't be null or empty")
var ErrEmailEmpty = errors.New("User email can't be null or empty")
var ErrInvalidEmailFormat = errors.New("Invalid email format")
var ErrPasswordEmpty = errors.New("User password can't be null or empty")
var ErrUserWithEmailExists = errors.New("User with the same email already exists")
