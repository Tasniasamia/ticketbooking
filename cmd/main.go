package main;

import (
	"ticketBooking/internal/config"
	"ticketBooking/internal/server"
	"ticketBooking/internal/database"


)




func main() {


 cfg:=config.LoadConfig(); // env variables available for OS
 db:=database.ConnectDatabase(cfg); // connect with database
 database.RunMigrations(db) // migrate all table


  server.Start(cfg,db);



 
}



