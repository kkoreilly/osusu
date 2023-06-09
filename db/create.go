package db

// Init creates all of the database tables if they do not exist.
func Init() error {
	err := CreateHstoreExtension()
	if err != nil {
		return err
	}
	err = CreateUsersTable()
	if err != nil {
		return err
	}
	err = CreateSessionsTable()
	if err != nil {
		return err
	}
	err = CreateGroupsTable()
	if err != nil {
		return err
	}
	err = CreateMealsTable()
	if err != nil {
		return err
	}
	err = CreateEntriesTable()
	if err != nil {
		return err
	}
	return nil
}

// CreateHstoreExtension creates the hstore extension in the database if it does not exist
func CreateHstoreExtension() error {
	statement := `CREATE EXTENSION IF NOT EXISTS hstore`
	_, err := db.Exec(statement)
	return err
}

// CreateUsersTable creates the users table in the database if it does not exist
func CreateUsersTable() error {
	statement := `
	CREATE TABLE IF NOT EXISTS public.users
	(
		id bigserial NOT NULL,
		username text NOT NULL,
		password text NOT NULL,
		name text NOT NULL,
		CONSTRAINT users_pkey PRIMARY KEY (id),
		CONSTRAINT users_username_key UNIQUE (username)
	)`
	_, err := db.Exec(statement)
	return err
}

// CreateSessionsTable creates the sessions table in the database if it does not eixst
func CreateSessionsTable() error {
	statement := `
	CREATE TABLE IF NOT EXISTS public.sessions
	(
		id text NOT NULL,
		user_id bigint NOT NULL,
		expires date NOT NULL,
		CONSTRAINT sessions_pkey PRIMARY KEY (id),
		CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id)
			REFERENCES public.users (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE
	)`
	_, err := db.Exec(statement)
	return err
}

// CreateGroupsTable creates the groups table in the database if it does not exist
func CreateGroupsTable() error {
	statement := `
	CREATE TABLE IF NOT EXISTS public.groups
	(
		id bigserial NOT NULL,
		owner bigint NOT NULL,
		code text NOT NULL,
		name text NOT NULL DEFAULT ''::text,
		members bigint[] NOT NULL,
		cuisines text[] NOT NULL DEFAULT '{}'::text[],
		CONSTRAINT group_pkey PRIMARY KEY (id),
		CONSTRAINT group_code_key UNIQUE (code),
		CONSTRAINT group_owner_fkey FOREIGN KEY (owner)
			REFERENCES public.users (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE RESTRICT
	)`
	_, err := db.Exec(statement)
	return err
}

// CreateMealsTable creates the meals table in the database if it does not exist
func CreateMealsTable() error {
	statement := `
	CREATE TABLE IF NOT EXISTS public.meals
	(
		id bigserial NOT NULL,
		group_id bigint NOT NULL,
		name text NOT NULL DEFAULT ''::text,
		description text NOT NULL DEFAULT ''::text,
		source text NOT NULL DEFAULT ''::text,
		image text NOT NULL DEFAULT ''::text,
		category text[] NOT NULL DEFAULT '{}'::text[],
		cuisine text[] NOT NULL DEFAULT '{}'::text[],
		CONSTRAINT meals_pkey PRIMARY KEY (id),
		CONSTRAINT meals_group_id_fkey FOREIGN KEY (group_id)
			REFERENCES public.groups (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE
	)`
	_, err := db.Exec(statement)
	return err
}

// CreateEntriesTable creates the entries table in the database if it does not exist
func CreateEntriesTable() error {
	statement := `
	CREATE TABLE IF NOT EXISTS public.entries
	(
		id bigserial NOT NULL,
		group_id bigint NOT NULL,
		meal_id bigint NOT NULL,
		entry_date date NOT NULL,
		category text NOT NULL DEFAULT ''::text,
		source text NOT NULL DEFAULT ''::text,
		cost hstore NOT NULL DEFAULT ''::hstore,
		effort hstore NOT NULL DEFAULT ''::hstore,
		healthiness hstore NOT NULL DEFAULT ''::hstore,
		taste hstore NOT NULL DEFAULT ''::hstore,
		CONSTRAINT entries_pkey PRIMARY KEY (id),
		CONSTRAINT entries_meal_id_fkey FOREIGN KEY (meal_id)
			REFERENCES public.meals (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE,
		CONSTRAINT entries_group_id_fkey FOREIGN KEY (group_id)
			REFERENCES public.groups (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE
	)`
	_, err := db.Exec(statement)
	return err
}
