package main

import (
	"log"

	"fmt"

	pb "github.com/arjunyel/student-info-api"
	"github.com/golang/protobuf/ptypes/empty"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStudentsClient(conn)

	stud1 := pb.Student{
		Id:    10,
		FName: "Test",
		LName: "Ing",
		Year:  4,
		Gpa:   10,
		Major: "art",
	}
	stud2 := pb.Student{
		Id:    120,
		FName: "Again",
		LName: "over",
		Year:  5,
		Gpa:   11,
		Major: "es",
	}
	stud1Request := pb.GetStudentRequest{
		Id: 10,
	}

	create1, err := c.CreateStudent(ctx, &stud1)
	if err != nil {
		log.Fatalf("Failed at create 1: %v", err)
	}
	fmt.Println(create1)
	create2, err := c.CreateStudent(ctx, &stud2)
	fmt.Println(create2)
	if err != nil {
		log.Fatalf("Failed at create 2: %v", err)
	}
	get1, err := c.GetStudent(ctx, &stud1Request)
	if err != nil {
		log.Fatalf("Failed at get 1: %v", err)
	}
	fmt.Println(get1)
	all, err := c.GetAllStudents(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("Failed at return all: %v", err)
	}
	fmt.Println(all)
}
