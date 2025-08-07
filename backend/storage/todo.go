package storage

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	todo "github.com/delta-team-uz/grpc-web-todo/todo_service_grpc"
	jsonitor "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

var json = jsonitor.ConfigCompatibleWithStandardLibrary

type TodoService interface {
	CreateTodo(ctx context.Context, req *todo.CreateTodoRequest) (*todo.CreateTodoResponse, error)
	GetAllTodo(ctx context.Context, req *emptypb.Empty) (*todo.GetAllResponse, error)
	DeleteTodo(ctx context.Context, req *todo.DeleteTodoRequest) (*emptypb.Empty, error)
	UpdateTodo(ctx context.Context, req *todo.UpdateTodoRequest) (*emptypb.Empty, error)
}

type todoService struct{}

func NewTodoService() TodoService {
	return &todoService{}
}

var dir string
var err error

func init() {
	dir, err = os.Getwd()
	if err != nil {
		slog.Error("error getting working directory", "error", err)
	}
	dir = dir + "/storage/todo.json"
}

func (t *todoService) CreateTodo(ctx context.Context, req *todo.CreateTodoRequest) (*todo.CreateTodoResponse, error) {
	var resp todo.CreateTodoResponse
	mo := protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Create(dir)
		os.WriteFile(dir, []byte("[]"), 0644)
	}

	bytes, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	var todos []*todo.Todo
	err = json.Unmarshal(bytes, &todos)
	if err != nil {
		slog.Error("error unmarshaling todo.json", "error", err)
		return nil, err
	}

	todo := &todo.Todo{
		Id:        strconv.Itoa(len(todos) + 1),
		Text:      req.GetText(),
		Completed: false,
	}

	todos = append(todos, todo)

	out := make([]jsonitor.RawMessage, 0, len(todos))
	for _, itm := range todos {
		b, merr := mo.Marshal(itm)
		if merr != nil {
			slog.Error("error marshaling todo.json", "error", merr)
			return nil, merr
		}
		out = append(out, jsonitor.RawMessage(b))
	}

	bytes, err = json.Marshal(out)
	if err != nil {
		slog.Error("error marshaling todo.json", "error", err)
		return nil, err
	}

	err = os.WriteFile(dir, bytes, 0644)
	if err != nil {
		slog.Error("error writing todo.json", "error", err)
		return nil, err
	}

	resp.Id = todo.GetId()

	return &resp, nil
}

func (t *todoService) DeleteTodo(ctx context.Context, req *todo.DeleteTodoRequest) (*emptypb.Empty, error) {
	bytes, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	var todos []*todo.Todo
	err = json.Unmarshal(bytes, &todos)
	if err != nil {
		slog.Error("error unmarshaling todo.json", "error", err)
		return nil, err
	}

	for i, t := range todos {
		if t.Id == req.Id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}

	bytes, err = json.Marshal(todos)
	if err != nil {
		slog.Error("error marshaling todo.json", "error", err)
		return nil, err
	}

	err = os.WriteFile(dir, bytes, 0644)
	if err != nil {
		slog.Error("error writing todo.json", "error", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (t *todoService) GetAllTodo(ctx context.Context, req *emptypb.Empty) (*todo.GetAllResponse, error) {
	bytes, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	var todos []*todo.Todo
	err = json.Unmarshal(bytes, &todos)
	if err != nil {
		slog.Error("error unmarshaling todo.json", "error", err)
		return nil, err
	}

	return &todo.GetAllResponse{
		Todo: todos,
	}, nil
}

func (t *todoService) UpdateTodo(ctx context.Context, req *todo.UpdateTodoRequest) (*emptypb.Empty, error) {
	bytes, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	var todos []*todo.Todo
	err = json.Unmarshal(bytes, &todos)
	if err != nil {
		slog.Error("error unmarshaling todo.json", "error", err)
		return nil, err
	}

	for _, t := range todos {
		if t.Id == req.Id {
			t.Text = req.Text
			t.Completed = req.Completed
		}
	}

	bytes, err = json.Marshal(todos)
	if err != nil {
		slog.Error("error marshaling todo.json", "error", err)
		return nil, err
	}

	err = os.WriteFile(dir, bytes, 0644)
	if err != nil {
		slog.Error("error writing todo.json", "error", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
