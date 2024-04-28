package shared

type GameStatus string

const (
	Status_New        GameStatus = "New"
	Status_Installing GameStatus = "installing"
	Status_Installed  GameStatus = "installed"
	Status_Error      GameStatus = "error"
)
