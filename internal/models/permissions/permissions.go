package permissions

type Permission string

const (
	All   Permission = "*"
	Read             = "r"
	Write            = "w"
)
