package manager

type Manager struct {
	ManagerID  string   `json:"manager_id" bson:"manager_id"`
	Gladiators []string `json:"gladiators" bson:"gladiators"`
}

func NewManager(managerID string) *Manager {
	m := &Manager{
		ManagerID: managerID,
	}

	return m
}
