package users

type Filter struct {
	FirstName string
	LastName string
}

func (f *Filter) getWhere() (where string, dest []interface{}) {
	if f.FirstName != "" {
		f.FirstName =  f.FirstName + "%"
		where += "first_name LIKE ?"
		dest = append(dest, f.FirstName)
	}

	if f.LastName != "" {
		if where != "" {
			where += " AND "
		}
		f.LastName =  f.LastName + "%"
		where += "last_name LIKE ?"
		dest = append(dest, f.LastName)
	}
	if where != "" {
		where = "WHERE " + where
	}

	return
}