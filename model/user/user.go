package user

import "os/user"

type User struct {
	Uid      string
	Username string
	Groups   []string
}

func New() *User {
	currentUser, err := user.Current()
	newUser := User{
		Uid:      currentUser.Uid,
		Username: currentUser.Username,
		Groups:   []string{},
	}
	if err != nil {
		// handle error
	}
	groupsId, err := currentUser.GroupIds()
	if err != nil {
		// handle error
	}

	for _, id := range groupsId {
		group, err := user.LookupGroupId(id)
		if err != nil {
			// handle error
		}
		newUser.Groups = append(newUser.Groups, group.Name)

	}
	return &newUser
}
