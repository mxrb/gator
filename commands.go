package main

import "errors"

type commandHandler func(*state, command) error

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]commandHandler
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
