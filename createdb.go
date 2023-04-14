package main

// CreateTablesDB creates all of the database tables if they do not exist.
func CreateTablesDB() error {
	err := CreateUsersTableDB()
	if err != nil {
		return err
	}
	err = CreateSessionsTableDB()
	if err != nil {
		return err
	}
	err = CreatePeopleTableDB()
	if err != nil {
		return err
	}
	err = CreateMealsTableDB()
	if err != nil {
		return err
	}
	err = CreateEntriesTableDB()
	if err != nil {
		return err
	}
	return nil
}

// CreateUsersTableDB creates the users table in the database if it does not exist
func CreateUsersTableDB() error {

}

// CreateSessionsTableDB creates the sessions table in the database if it does not eixst
func CreateSessionsTableDB() error {

}

// CreatePeopleTableDB creates the people table in the database if it does not exist
func CreatePeopleTableDB() error {

}

// CreateMealsTableDB creates the meals table in the database if it does not exist
func CreateMealsTableDB() error {

}

// CreateEntriesTableDB creates the entries table in the database if it does not exist
func CreateEntriesTableDB() error {

}
