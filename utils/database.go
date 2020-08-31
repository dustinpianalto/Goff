package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	Database *sql.DB
)

func ConnectDatabase(dbConnString string) {
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		panic(fmt.Sprintf("Can't connect to the database. %v", err))
	} else {
		fmt.Println("Database Connected.")
	}
	Database = db
}

func InitializeDatabase() {
	_, err := Database.Query("CREATE TABLE IF NOT EXISTS users(" +
		"id varchar(30) primary key," +
		"banned bool not null default false," +
		"logging bool not null default true," +
		"steam_id varchar(30) NOT NULL DEFAULT ''," +
		"is_active bool not null default true," +
		"is_staff bool not null default false," +
		"is_admin bool not null default false" +
		")")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("CREATE TABLE IF NOT EXISTS guilds(" +
		"id varchar(30) primary key," +
		"welcome_message varchar(1000) NOT NULL DEFAULT ''," +
		"goodbye_message varchar(1000) NOT NULL DEFAULT ''," +
		"logging_channel varchar(30) NOT NULL DEFAULT ''," +
		"welcome_channel varchar(30) NOT NULL DEFAULT ''" +
		")")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("CREATE TABLE IF NOT EXISTS prefixes(" +
		"id serial primary key," +
		"prefix varchar(10) not null unique default 'Go.'" +
		")")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("CREATE TABLE IF NOT EXISTS tags(" +
		"id serial primary key," +
		"tag varchar(100) not null unique," +
		"content varchar(1000) not null," +
		"creator varchar(30) not null references users(id)," +
		"creation_time timestamp not null default NOW()," +
		"guild_id varchar(30) not null references guilds(id)" +
		")")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("CREATE TABLE IF NOT EXISTS x_users_guilds(" +
		"guild_id varchar(30) not null references guilds(id)," +
		"user_id varchar(30) not null references users(id)" +
		")")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("CREATE TABLE IF NOT EXISTS x_guilds_prefixes(" +
		"guild_id varchar(30) not null references guilds(id)," +
		"prefix_id int not null references prefixes(id)" +
		")")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("CREATE TABLE IF NOT EXISTS tasks(" +
		"id serial primary key," +
		"type varchar(10) not null," +
		"content text not null," +
		"guild_id varchar(30) not null references guilds(id)," +
		"channel_id varchar(30) not null," +
		"user_id varchar(30) not null," +
		"creation_time timestamp not null default NOW()," +
		"trigger_time timestamp not null," +
		"completed bool not null default false," +
		"processing bool default false)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query(`CREATE TABLE IF NOT EXISTS postfixes(
		id serial primary key,
		name varchar(100) not null,
		time timestamp not null default NOW())`)
	if err != nil {
		log.Println(err)
	}
	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS puzzles(
		id serial primary key,
		text text not null,
		time timestamp not null
		)`)
	if err != nil {
		log.Println(err)
	}
	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS x_guilds_puzzles(
		id serial primary key,
		guild_id varchar(30) not null references guilds(id),
		puzzle_id int not null references puzzles(id),
		message_id varchar(30) not null
		)`)
	RunPostfixes()
}

func LoadTestData() {
	_, err := Database.Query("INSERT INTO users (id, banned, logging, steam_id, is_active, is_staff, is_admin) values " +
		"('351794468870946827', false, true, '76561198024193239', true, true, true)," +
		"('692908139506434065', false, true, '', true, false, false)," +
		"('396588996706304010', false, true, '', true, true, false)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("INSERT INTO guilds (id, welcome_message, goodbye_message) VALUES " +
		"('265828729970753537', 'Hey there is someone new here.', 'Well fine then... Just leave without saying goodbye')")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("INSERT INTO prefixes (prefix) VALUES ('Godev.'), ('godev.'), ('godev,')")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("INSERT INTO x_users_guilds (guild_id, user_id) VALUES " +
		"('265828729970753537', '351794468870946827')," +
		"('265828729970753537', '692908139506434065')," +
		"('265828729970753537', '396588996706304010')")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("INSERT INTO x_guilds_prefixes (guild_id, prefix_id) VALUES " +
		"('265828729970753537', 1)," +
		"('265828729970753537', 2)," +
		"('265828729970753537', 3)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Database.Query("INSERT INTO tags (tag, content, creator, guild_id) VALUES " +
		"('test', 'This is a test of the tag system', '351794468870946827', '265828729970753537')")
	if err != nil {
		fmt.Println(err)
	}
}
