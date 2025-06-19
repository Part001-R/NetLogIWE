package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Part001-R/netlogiwe/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/joho/godotenv"

	db "github.com/Part001-R/netlogiwe/pkg/db"
)

type server struct {
	pb.UnimplementedIweServer
	db db.ActionsDB
}

func main() {

	// Preparatory actions
	objDB, closeDb, err := preparAct()
	if err != nil {
		log.Fatalf("fault preparatory actions: %v", err)
	}
	defer func() {
		err := closeDb()
		if err != nil {
			fmt.Println("fault close DB")
		}
	}()

	// gRPCS
	srvImpl := &server{
		db: objDB,
	}
	err = startUpServer(srvImpl)
	if err != nil {
		log.Fatalf("fault start up IWE server: %v", err)
	}
}

// Handler
func (s *server) SaveMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {

	var msg = db.MessageT{}

	msg.TypeMessage = req.GetTypeMessage()
	msg.NameProject = req.GetNameProject()
	msg.LocationEvent = req.GetLocationEvent()
	msg.BodyMessage = req.GetBodyMessage()

	if msg.TypeMessage == "T" {
		return &pb.MessageResponse{Status: "Ok"}, nil
	}

	err := s.db.SavingMessage(msg)

	//err := db.StoreMessage(s.db, msg)
	if err != nil {
		fmt.Printf("error: {%v}\n", err)
		return nil, err
	}

	return &pb.MessageResponse{Status: "Ok"}, nil
}

// preparatory actions. Returns: db pointer, function close db connect, error
func preparAct() (db.ActionsDB, func() error, error) {

	// ENV
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Fault read env file")
	}

	// DB
	ptrDb, close, err := db.ConDb(os.Getenv("DB_TYPE"), os.Getenv("DB_NAME"))
	if err != nil {
		return nil, nil, fmt.Errorf("fault connect DB: %v", err)
	}
	objDB, err := db.RepoDB(ptrDb)
	if err != nil {
		log.Fatalf("an error create object instance: '%v'", err)
	}

	// Tables
	err = objDB.Tables()
	if err != nil {
		return nil, close, fmt.Errorf("fault create tables: %v", err)
	}
	return objDB, close, nil
}

// Start up IWE server. Return error.
func startUpServer(s *server) error {

	creds, err := credentials.NewServerTLSFromFile(os.Getenv("PATH_PUBLIC_KEY"), os.Getenv("PATH_PRIVATE_KEY"))
	if err != nil {
		return fmt.Errorf("fault read sertificats: %v", err)
	}

	ipAndPort := os.Getenv("PORT")
	listener, err := net.Listen("tcp", ipAndPort)
	if err != nil {
		return fmt.Errorf("fault create listener tcp port %s: %v", ipAndPort, err)
	}

	srv := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterIweServer(srv, s)
	log.Println("Start up IWE server:", ipAndPort)

	err = srv.Serve(listener)
	if err != nil {
		return fmt.Errorf("fault start up IWE server: %v ", err)
	}

	return errors.New("plug error")
}
