package module

type Config interface {
	Initialize(context *Context, name string)
	Validate()
}

type DefaultConfig struct {
}

func (c *DefaultConfig) Validate() {

}
