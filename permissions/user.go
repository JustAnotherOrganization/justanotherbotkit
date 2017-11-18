package permissions

type (
	// User is a standard chat user.
	User struct {
		ID    string
		perms []string
		// TODO: implement groups

		pm *Manager
	}
)

// GetPerm returns true if the user has the selected perm.
func (u *User) GetPerm(perm string) (bool, error) {
	perms, err := u.GetPerms()
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p == perm {
			return true, nil
		}
	}

	return false, nil
}

// GetPerms returns all the perms for the user.
func (u *User) GetPerms() ([]string, error) {
	return u.pm.db.GetPerms(u.ID)
}

// AddPerms adds perms to the user.
func (u *User) AddPerms(perms ...string) error {
	_perms, err := u.pm.db.GetPerms(u.ID)
	if err != nil {
		return err
	}

	changed := false
	for _, p := range perms {
		exists := false

		for _, _p := range _perms {
			if p == _p {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		_perms = append(_perms, p)
		changed = true
	}

	if !changed {
		return nil
	}

	return u.pm.db.SetPerms(u.ID, _perms...)
}

// DelPerms deletes perms from the user.
func (u *User) DelPerms(perms ...string) error {
	_perms, err := u.pm.db.GetPerms(u.ID)
	if err != nil {
		return err
	}

	var (
		final []string
		found bool
	)
	for _, _p := range _perms {
		for _, p := range perms {
			if _p == p {
				found = true
			}
		}

		if !found {
			final = append(final, _p)
		}
	}

	return u.pm.db.SetPerms(u.ID, final...)
}
