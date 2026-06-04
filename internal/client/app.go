package client

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
	"github.com/LemuriiL/GophKeeperPassword/internal/secure"
)

// App хранит клиентское приложение
type App struct {
	cfg Config
	api *APIClient
}

// NewApp создает приложение клиента
func NewApp(configShort string, configLong string) (*App, error) {
	cfg, err := LoadConfig(configShort, configLong)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg: cfg,
		api: NewAPIClient(cfg.ServerURL),
	}, nil
}

// Run запускает CLI команду
func (a *App) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("commands: register, login, add, get, list, update, delete")
	}

	switch args[0] {
	case "register":
		return a.runRegister(args[1:])
	case "login":
		return a.runLogin(args[1:])
	case "add":
		return a.runAdd(args[1:])
	case "get":
		return a.runGet(args[1:])
	case "list":
		return a.runList(args[1:])
	case "update":
		return a.runUpdate(args[1:])
	case "delete":
		return a.runDelete(args[1:])
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

	token, _, err := a.api.Register(*login, *password)
	if err != nil {
		return err
	}

	return SaveSession(a.cfg.SessionFile, Session{
		Login: *login,
		Token: token,
	})
}

func (a *App) runLogin(args []string) error {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)

	login := fs.String("login", "", "login")
	password := fs.String("password", "", "password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	token, err := a.api.Login(*login, *password)
	if err != nil {
		return err
	}

	session, loadErr := LoadSession(a.cfg.SessionFile)
	if loadErr != nil {
		session = Session{}
	}

	session.Login = *login
	session.Token = token

	return SaveSession(a.cfg.SessionFile, session)
}

func (a *App) runAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)

	kind := fs.String("type", model.TypeText, "item type")
	title := fs.String("title", "", "title")
	meta := fs.String("meta", "", "meta")
	value := fs.String("value", "", "value json")
	password := fs.String("master-password", "", "master password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	session, err := LoadSession(a.cfg.SessionFile)
	if err != nil {
		return err
	}

	itemSalt, err := secure.NewSalt()
	if err != nil {
		return err
	}

	ciphertext, nonce, err := EncryptPayload(*password, itemSalt, *value)
	if err != nil {
		return err
	}

	item, err := a.api.SaveItem(session.Token, dto.UpsertItemRequest{
		Type:       *kind,
		Title:      *title,
		Meta:       *meta,
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Salt:       itemSalt,
	})
	if err != nil {
		return err
	}

	fmt.Println(item.ID)
	return nil
}

func (a *App) runGet(args []string) error {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)

	id := fs.String("id", "", "item id")
	password := fs.String("master-password", "", "master password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	session, err := LoadSession(a.cfg.SessionFile)
	if err != nil {
		return err
	}

	item, err := a.api.GetItem(session.Token, *id)
	if err != nil {
		return err
	}

	printItem(item, session.Salt, *password)
	return nil
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
		printItem(item, session.Salt, *password)
	}

	return nil
}

func (a *App) runUpdate(args []string) error {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)

	id := fs.String("id", "", "item id")
	kind := fs.String("type", model.TypeText, "item type")
	title := fs.String("title", "", "title")
	meta := fs.String("meta", "", "meta")
	value := fs.String("value", "", "value json")
	password := fs.String("master-password", "", "master password")

	if err := fs.Parse(args); err != nil {
		return err
	}

	session, err := LoadSession(a.cfg.SessionFile)
	if err != nil {
		return err
	}

	itemSalt, err := secure.NewSalt()
	if err != nil {
		return err
	}

	ciphertext, nonce, err := EncryptPayload(*password, itemSalt, *value)
	if err != nil {
		return err
	}

	_, err = a.api.UpdateItem(session.Token, *id, dto.UpsertItemRequest{
		Type:       *kind,
		Title:      *title,
		Meta:       *meta,
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Salt:       itemSalt,
	})

	return err
}

func (a *App) runDelete(args []string) error {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)

	id := fs.String("id", "", "item id")

	if err := fs.Parse(args); err != nil {
		return err
	}

	session, err := LoadSession(a.cfg.SessionFile)
	if err != nil {
		return err
	}

	return a.api.DeleteItem(session.Token, *id)
}

func printItem(item dto.ItemResponse, sessionSalt string, password string) {
	salt := item.Salt
	if strings.TrimSpace(salt) == "" {
		salt = sessionSalt
	}

	plain, err := DecryptPayload(password, salt, item.Ciphertext, item.Nonce)
	if err != nil {
		plain = ""
	}

	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("ID:", item.ID)
	fmt.Println("Type:", item.Type)
	fmt.Println("Title:", item.Title)
	fmt.Println("Meta:", item.Meta)
	fmt.Println("Value:", plain)
	fmt.Println("Updated:", item.UpdatedAt)
}
