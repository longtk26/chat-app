package presenters

import "github.com/longtk26/chat-app/internal/usecases"

type SocketPresenter struct {
	us usecases.ISocketUsecase
}

func NewSocketPresenter(us usecases.ISocketUsecase) *SocketPresenter {
	return &SocketPresenter{
		us: us,
	}
}

func (p *SocketPresenter) HandleConnect() {
	p.us.HandleConnect()
}
