package databases

type User struct {
	ID        string
	Firstname string
	Lastname  string
	Email     string
	Password  string
}

var usersDB []User

type GlobalError struct{}

func GetUserByEmail(email string) (*GlobalError, *User) {
	for _, user := range usersDB {
		if user.Email == email {
			return nil, &user
		}
	}

	return &GlobalError{}, nil
}

func GetUserByID(ID string) (*GlobalError, *User) {
	for _, user := range usersDB {
		if user.ID == ID {
			return nil, &user
		}
	}

	return &GlobalError{}, nil
}

func CreateUser(user User) *GlobalError {
	err, existingUser := GetUserByID(user.ID)

	if err != nil || existingUser != nil {
		return &GlobalError{}
	}

	usersDB = append(usersDB, user)
	return nil
}

func GetAllUsers() []User {
	return usersDB
}
