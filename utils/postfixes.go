package utils

import "log"

type postfix struct {
	Name   string
	Invoke func(bool) error
}

var postfixes = []postfix{
	postfix{
		Name:   "1_Update_Guild_for_Puzzle",
		Invoke: updateGuildForPuzzle,
	},
	postfix{
		Name:   "1_Update_X_Guild_Prefixes_to_add_ID",
		Invoke: updateXGuildPrefixesToAddID,
	},
	postfix{
		Name:   "1_Update_Tags_Content_Length",
		Invoke: updateTagsContentLength,
	},
}

func RunPostfixes() {
	for _, postfix := range postfixes {
		queryString := "SELECT * from postfixes where name = $1"
		rows, err := Database.Query(queryString, postfix.Name)
		if err != nil {
			log.Println(err)
			continue
		}
		if rows.Next() {
			continue
		} else {
			err := postfix.Invoke(false)
			if err != nil {
				continue
			}
			_, err = Database.Exec("INSERT INTO postfixes (name) VALUES ($1)", postfix.Name)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func updateGuildForPuzzle(revert bool) error {
	var queryString string
	if !revert {
		queryString = `ALTER TABLE guilds
		ADD COLUMN puzzle_channel varchar(30) not null default ''`
	} else {
		queryString = `ALTER TABLE guilds
		DROP COLUMN puzzleChat`
	}
	_, err := Database.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func updateXGuildPrefixesToAddID(revert bool) error {
	var queryString string
	if !revert {
		queryString = `ALTER TABLE x_guilds_prefixes
		ADD COLUMN id serial primary key`
	} else {
		queryString = `ALTER TABLE x_guilds_prefixes
		DROP COLUMN id`
	}
	_, err := Database.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func updateTagsContentLength(revert bool) error {
	var queryString string
	if !revert {
		queryString = `ALTER TABLE tags
		ALTER COLUMN content TYPE varchar(2000)`
	} else {
		queryString = `ALTER TABLE tags
		ALTER COLUMN content TYPE varchar(1000)`
	}
	_, err := Database.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
