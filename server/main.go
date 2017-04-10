package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"

	context "golang.org/x/net/context"

	pb "github.com/arjunyel/student-info-api"
	"github.com/arjunyel/student-info-api/database"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8000"
)

type server struct{}
type student struct {
	id    int32
	fName string
	lName string
	year  int32
	gpa   int32
	major string
}

func (s *server) CreateStudent(ctx context.Context, in *pb.Student) (*pb.Student, error) {
	kind := "Student"
	id := strconv.Itoa(int(in.Id))
	studentKey := datastore.NameKey(kind, id, nil)
	if _, err := database.DB.DB.Put(ctx, studentKey, in); err != nil {
		return nil, err
	}
	fmt.Printf("Saved %v\n", studentKey)
	return in, nil
}

func (s *server) GetAllStudents(ctx context.Context, in *empty.Empty) (*pb.AllStudents, error) {
	var students []*pb.Student
	fmt.Println("Get all")
	query := datastore.NewQuery("Student")
	_, err := database.DB.DB.GetAll(ctx, query, &students)
	if err != nil {
		return nil, err
	}
	return &pb.AllStudents{Students: students}, nil
}

func (s *server) GetStudent(ctx context.Context, in *pb.GetStudentRequest) (*pb.Student, error) {
	var student pb.Student
	id := strconv.Itoa(int(in.Id))
	key := datastore.NameKey("Student", id, nil)
	err := database.DB.DB.Get(ctx, key, &student)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(student)
	return &student, nil
}

func main() {
	projID := os.Getenv("DATASTORE_PROJECT_ID")
	if projID == "" {
		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
	}
	// [START build_service]
	ctx := context.Background()
	var err error
	database.DB.DB, err = datastore.NewClient(ctx, projID)
	// [END build_service]
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	// Initializes the gRPC server.
	s := grpc.NewServer()

	// Register the server with gRPC.
	pb.RegisterStudentsServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
