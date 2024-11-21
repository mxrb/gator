package main

import (
	"context"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}
	userName := cmd.Args[0]
	dbUser, err := s.db.GetUserByName(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	if err := s.cfg.SetUser(dbUser.Name); err != nil {
		return err
	}
	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}
	userName := cmd.Args[0]
	dbUser, err := s.db.CreateUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	if err := s.cfg.SetUser(userName); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("User created successfully: %s", dbUser)
	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset users: %w", err)
	}
	fmt.Println("Users were reset successfully!")
	return nil
}

func handlerUsers(s *state, _ command) error {
	dbUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetch users: %w", err)
	}

	for _, dbUser := range dbUsers {
		if dbUser.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", dbUser.Name)
		} else {
			fmt.Printf("* %s\n", dbUser.Name)
		}
	}
	return nil
}
