package permissions

func isValidPermissions(permissions uint8) bool {
	switch permissions {
	case PermissionNone:
	case PermissionRead:
	case PermissionWrite:
	case PermissionReadWrite:
	default:
		return false
	}

	return true
}
