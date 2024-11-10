package session

type ISessionRepository interface {
	Create(uint64, string) error
	GetUserIdBySession(string) (uint64, error)
	Delete(string) error
}

type ISessionUsecase interface {
	// GetUserFromSession(session string) (uint64, error)
}
