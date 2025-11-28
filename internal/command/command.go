package command

// Command is the base interface for all CLI commands.
type Command interface {
	Name() string
	Description() string
	Run(args []string) error
}

// Registry holds registered commands.
type Registry struct {
	cmds map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{cmds: make(map[string]Command)}
}

func (r *Registry) Register(c Command) {
	r.cmds[c.Name()] = c
}

func (r *Registry) Get(name string) (Command, bool) {
	c, ok := r.cmds[name]
	return c, ok
}

func (r *Registry) All() []Command {
	out := make([]Command, 0, len(r.cmds))
	for _, c := range r.cmds {
		out = append(out, c)
	}
	return out
}
