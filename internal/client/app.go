package client

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
)

type App struct {
	cfg Config
	api *APIClient
}

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

	token, salt, err := a.api.Register(*login, *password)
	if err != nil {
		return err
	}

	if strings.TrimSpace(salt) == "" {
		return errors.New("server returned empty salt")
	}

	return SaveSession(a.cfg.SessionFile, Session{
		Login: *login,
		Token: token,
		Salt:  salt,
	})
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

	if strings.TrimSpace(salt) == "" {
		return errors.New("server returned empty salt")
	}

	return SaveSession(a.cfg.SessionFile, Session{
		Login: *login,
		Token: token,
		Salt:  salt,
	})
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

	if strings.TrimSpace(session.Salt) == "" {
		return errors.New("empty session salt, run login again")
	}

	ciphertext, nonce, err := EncryptPayload(*password, session.Salt, *value)
	if err != nil {
		return err
	}

	check, err := DecryptPayload(*password, session.Salt, ciphertext, nonce)
	if err != nil {
		return fmt.Errorf("local encryption self-check failed: %w", err)
	}

	if check != *value {
		return errors.New("local encryption self-check failed: plaintext mismatch")
	}

	item, err := a.api.SaveItem(session.Token, dto.UpsertItemRequest{
		Type:       *kind,
		Title:      *title,
		Meta:       *meta,
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Salt:       session.Salt,
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

	if strings.TrimSpace(session.Salt) == "" {
		return errors.New("empty session salt, run login again")
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

	if strings.TrimSpace(session.Salt) == "" {
		return errors.New("empty session salt, run login again")
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

	if strings.TrimSpace(session.Salt) == "" {
		return errors.New("empty session salt, run login again")
	}

	ciphertext, nonce, err := EncryptPayload(*password, session.Salt, *value)
	if err != nil {
		return err
	}

	check, err := DecryptPayload(*password, session.Salt, ciphertext, nonce)
	if err != nil {
		return fmt.Errorf("local encryption self-check failed: %w", err)
	}

	if check != *value {
		return errors.New("local encryption self-check failed: plaintext mismatch")
	}

	_, err = a.api.UpdateItem(session.Token, *id, dto.UpsertItemRequest{
		Type:       *kind,
		Title:      *title,
		Meta:       *meta,
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Salt:       session.Salt,
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
	plain, err := DecryptPayload(password, sessionSalt, item.Ciphertext, item.Nonce)
	if err != nil {
		plain = fmt.Sprintf("<decrypt error: %v>", err)
	}

	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("ID:", item.ID)
	fmt.Println("Type:", item.Type)
	fmt.Println("Title:", item.Title)
	fmt.Println("Meta:", item.Meta)
	fmt.Println("Value:", plain)
	fmt.Println("Updated:", item.UpdatedAt)
}
