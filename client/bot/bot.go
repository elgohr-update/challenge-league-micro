// Package bot is a Hubot style bot that sits a microservice environment
package bot

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"

	"github.com/micro/go-micro/v2/agent/command"
	"github.com/micro/go-micro/v2/agent/input"
	log "github.com/micro/go-micro/v2/logger"

	// inputs
	nakamaCommands "github.com/challenge-league/nakama-go/commands"
	nakama "github.com/challenge-league/nakama-go/context"
	_ "github.com/micro/go-micro/v2/agent/input/discord"
	_ "github.com/micro/go-micro/v2/agent/input/slack"
	_ "github.com/micro/go-micro/v2/agent/input/telegram"
)

type bot struct {
	exit chan bool
	ctx  *cli.Context

	sync.RWMutex
	inputs   map[string]input.Input
	commands map[string]command.Command
}

var (
	// Default server name
	Name = "go.micro.bot"
	// Namespace for commands
	Namespace = "go.micro.bot"
	// map pattern:command
)

func newBot(ctx *cli.Context, inputs map[string]input.Input, commands map[string]command.Command) *bot {
	return &bot{
		ctx:      ctx,
		exit:     make(chan bool),
		commands: commands,
		inputs:   inputs,
	}
}

func (b *bot) loop(io input.Input) {
	log.Infof("[loop] starting %s", io.String())

	for {
		select {
		case <-b.exit:
			log.Infof("[loop] exiting %s", io.String())
			return
		default:
			if err := b.run(io); err != nil {
				log.Error("[loop] error %v", err)
				time.Sleep(time.Second)
			}
		}
	}
}

func (b *bot) process(c input.Conn, ev input.Event) error {
	args := strings.Split(string(ev.Data), " ")
	if len(args) == 0 {
		return nil
	}

	b.RLock()
	defer b.RUnlock()

	nakamaCtx, err := nakama.NewCustomAuthenticatedDiscordAPIClient(ev.DiscordMsg)
	if err != nil {
		log.Error(err)
		c.Send(&input.Event{
			Meta: ev.Meta,
			From: ev.To,
			To:   ev.From,
			Type: input.TextEvent,
			Data: []byte(err.Error()),
		})
		return nil
	}
	defer nakamaCtx.Conn.Close()
	log.Infof("nakamaCtx %v", nakamaCtx)

	cmdBuilder := nakamaCommands.NewCommandsBuilder()
	cmdBuilder.SetContext(nakamaCtx)
	log.Infof("args %+v\n", args)
	output, err := nakamaCommands.ExecuteCommandC(cmdBuilder, args...)
	if err != nil {
		fmt.Printf("Error: %+v", err)
	}

	return c.Send(&input.Event{
		Meta: ev.Meta,
		From: ev.To,
		To:   ev.From,
		Type: input.TextEvent,
		Data: []byte(output),
	})
}

func (b *bot) run(io input.Input) error {
	log.Infof("[loop] connecting to %s", io.String())

	c, err := io.Stream()
	if err != nil {
		return err
	}

	for {
		select {
		case <-b.exit:
			log.Infof("[loop] closing %s", io.String())
			return c.Close()
		default:
			var recvEv input.Event
			// receive input
			if err := c.Recv(&recvEv); err != nil {
				return err
			}

			// only process TextEvent
			if recvEv.Type != input.TextEvent {
				continue
			}

			if len(recvEv.Data) == 0 {
				continue
			}

			if err := b.process(c, recvEv); err != nil {
				return err
			}
		}
	}
}

func (b *bot) start() error {
	log.Info("starting")

	// Start inputs
	for _, io := range b.inputs {
		log.Infof("starting input %s", io.String())

		if err := io.Init(b.ctx); err != nil {
			return err
		}

		if err := io.Start(); err != nil {
			return err
		}

		go b.loop(io)
	}

	return nil
}

func (b *bot) stop() error {
	log.Info("stopping")
	close(b.exit)

	// Stop inputs
	for _, io := range b.inputs {
		log.Infof("stopping input %s", io.String())
		if err := io.Stop(); err != nil {
			log.Errorf("%v", err)
		}
	}

	return nil
}

func run(ctx *cli.Context) error {
	log.Init(log.WithFields(map[string]interface{}{"service": "bot"}))

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// Parse inputs
	var inputs []string

	parts := strings.Split(ctx.String("inputs"), ",")
	for _, p := range parts {
		if len(p) > 0 {
			inputs = append(inputs, p)
		}
	}

	ios := make(map[string]input.Input)
	cmds := make(map[string]command.Command)

	// Parse inputs
	for _, io := range inputs {
		if len(io) == 0 {
			continue
		}
		i, ok := input.Inputs[io]
		if !ok {
			log.Errorf("input %s not found\n", i)
			os.Exit(1)
		}
		ios[io] = i
	}

	// setup service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(
			time.Duration(ctx.Int("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.Int("register_interval"))*time.Second,
		),
	)

	// Start bot
	b := newBot(ctx, ios, cmds)

	if err := b.start(); err != nil {
		log.Errorf("error starting bot %v", err)
		os.Exit(1)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop bot
	if err := b.stop(); err != nil {
		log.Errorf("error stopping bot %v", err)
	}

	return nil
}

func Commands() []*cli.Command {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "inputs",
			Usage:   "Inputs to load on startup",
			EnvVars: []string{"MICRO_BOT_INPUTS"},
		},
		&cli.StringFlag{
			Name:    "namespace",
			Usage:   "Set the namespace used by the bot to find commands e.g. com.example.bot",
			EnvVars: []string{"MICRO_BOT_NAMESPACE"},
		},
	}

	// setup input flags
	for _, input := range input.Inputs {
		flags = append(flags, input.Flags()...)
	}

	command := &cli.Command{
		Name:   "bot",
		Usage:  "Run the chatops bot",
		Flags:  flags,
		Action: run,
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []*cli.Command{command}
}
