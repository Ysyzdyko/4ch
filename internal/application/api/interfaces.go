package api

type UserService interface {
	GenUserID(userCount, maxID int) (int, error)
}

// type PostService interface {
// }
