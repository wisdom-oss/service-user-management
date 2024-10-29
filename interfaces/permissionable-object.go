package interfaces

type PermissionableObject interface {
	GetID() string
	Permissions() map[string][]string
	IsAdministrator() bool
}
