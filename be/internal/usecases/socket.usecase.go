package usecases

import "fmt"

type ISocketUsecase interface {
	HandleConnect()
}

type SocketUsecase struct{}

var _ ISocketUsecase = &SocketUsecase{}

func NewSocketUsecase() ISocketUsecase {
	return &SocketUsecase{}
}

func (s *SocketUsecase) HandleConnect() {
	fmt.Println("Handling socket connection...")
}
