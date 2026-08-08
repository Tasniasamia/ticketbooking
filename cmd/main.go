package main;

import (
	"ticketBooking/internal/config"
	"ticketBooking/internal/server"
	"ticketBooking/internal/database"


)




func main() {


 cfg:=config.LoadConfig();
 db:=database.ConnectDatabase(cfg);
 database.RunMigrations(db)


  server.Start(cfg,db);



 
}



