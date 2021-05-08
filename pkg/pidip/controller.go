package pidip

type Task interface {
	work()
}

type PIController struct {
	Kp float64
	Ki float64
	i float64
}

func NewPIController(kp, ki float64) PIController {
	return PIController{kp, ki, 0.0}
}

func (c *PIController) work(e float64) float64 {
	c.i += e

	return c.Kp * e +  c.Ki * c.i
}
