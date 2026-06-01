package client

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
)

type App struct {
	cfg Config
	api *APIClient
}

func NewApp() (*App, error) {
	cfg, err := LoadConfig("", "")
	if err != nil {
		return nil, err
	}

	return &App{
		cfg: cfg,
		api: NewAPIClient(cfg.ServerURL),
	}, nil
}

func (a *App) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("commands: register, login, add, list")
	}

	switch args[0] {
	case "register":
		return a.runRegister(args[1:])
	case "login":
		return a.runLogin(args[1:])
	case "add":
		return a.runAdd(args[1:])
	case "list":
		return a.runList(args[1:])
	default:
		return errors.New("unknown command")
	}
}

func (a *App) runRegister(args []string) error {
	fs := flag.NewFlagSet("register", flag.ContinueOnError)

	login := fs.String("login", "", "login")
	password := fs.String("password", "", "password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	return a.api.Register(*login, *password)
}

func (a *App) runLogin(args []string) error {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)

	login := fs.String("login", "", "login")
	password := fs.String("password", "", "password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	token, salt, err := a.api.Login(*login, *password)
	if err != nil {
		return err
	}

	return SaveSession(a.cfg.SessionFile, Session{
		Login: *login,
		Token: token,
		Salt:  salt,
	})
}

func (a *App) runAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)

	kind := fs.String("type", "text", "item type")
	title := fs.String("title", "", "title")
	meta := fs.String("meta", "", "meta")
	value := fs.String("value", "", "value")
	password := fs.String("master-password", "", "master password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	session, err := LoadSession(a.cfg.SessionFile)
	if err != nil {
		return err
	}

	ciphertext, nonce, err := EncryptPayload(*password, session.Salt, *value)
	if err != nil {
		return err
	}

	return a.api.SaveItem(session.Token, dto.UpsertItemRequest{
		Type:       *kind,
		Title:      *title,
		Meta:       *meta,
		Ciphertext: ciphertext,
		Nonce:      nonce,
	})
}

func (a *App) runList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)

	password := fs.String("master-password", "", "master password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	session, err := LoadSession(a.cfg.SessionFile)
	if err != nil {
		return err
	}

	items, err := a.api.ListItems(session.Token)
	if err != nil {
		return err
	}

	for _, item := range items {
		plain, err := DecryptPayload(*password, session.Salt, item.Ciphertext, item.Nonce)
		if err != nil {
			plain = "<decrypt error>"
		}

		fmt.Println(strings.Repeat("-", 40))
		fmt.Println("ID:", item.ID)
		fmt.Println("Type:", item.Type)
		fmt.Println("Title:", item.Title)
		fmt.Println("Meta:", item.Meta)
		fmt.Println("Value:", plain)
		fmt.Println("Updated:", item.UpdatedAt)
	}

	return nil
}
