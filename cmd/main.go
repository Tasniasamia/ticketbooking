package main

import (
	"ticketBooking/internal/config"
	"ticketBooking/internal/server"

)





func main() {

 cfg:=config.LoadConfig();
 db:=config.ConnectDatabase(cfg);



  server.Start(cfg,db);



 
}



