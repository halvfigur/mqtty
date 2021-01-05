package control

const startPageLabel = "startpage"

type (
	ConnectFunc         func(host string, port int, username, password string)
	StartPageController struct {
		ctrl Control
	}
)

func NewStartPageController(ctrl Control) *StartPageController {
	return &StartPageController{
		ctrl: ctrl,
	}
}

func (c *StartPageController) OnConnect(host string, port int, username, password string) {
	c.ctrl.Connect(host, port, username, password)
}

func (c *StartPageController) OnExit() {
	c.ctrl.OnStop()
}
