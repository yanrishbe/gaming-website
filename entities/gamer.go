package entities

type User struct {
	Id      int    `json:"id, omitempty"`
	Name    string `json:"name"`
	Balance *int    `json:"balance, omitempty"`
	//Points *int `json:"points, omitempty"`
}

var Users []User

func CreateUser(user User) User {
	user.Id = len(Users) + 1
	*user.Balance -= 300
	Users = append(Users, user)
	return user
}

func RemoveUser(id int) {
	for i := id; i < len(Users); i++ {
		Users[id].Id = i
	}
	Users = append(Users[:id-1], Users[id:]...)
}
