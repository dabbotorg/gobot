package main

import (
	"strings"

	"github.com/dabbotorg/gobot/commands"

	"github.com/foxbot/discordgo"
)

var cmds = commands.Commands()
var mentionPrefix string // will be set by onReady

func onMessage(s *discordgo.Session, e *discordgo.MessageCreate) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				errors <- e
			} else {
				logger.Println("recovered from a fatal non-error:", err)
			}
		}
	}()

	if e.Author.Bot {
		return
	}

	err := command(s, e)
	if err != nil {
		errors <- err
	}
}

func command(s *discordgo.Session, e *discordgo.MessageCreate) error {
	// TODO: temp, use guild prefix
	hasPrefix, argPos := checkPrefix(s, e)
	if !hasPrefix {
		return nil
	}
	c := e.Content[argPos:]
	p := strings.Split(c, " ")
	if len(p) < 1 {
		return nil
	}
	name := p[0]

	var args []string
	if len(p) < 2 {
		args = []string{}
	} else {
		args = p[1:]
	}

	cmd, ok := cmds[name]
	if !ok {
		return nil
	}

	ctx := &commands.Context{
		Args:     args,
		Config:   &conf,
		Event:    e,
		Lavalink: lavalink,
		Owo:      owoclient,
		Redis:    rdis,
		Session:  s,
	}

	r := cmd.Method(ctx)
	err := r.Act(ctx)

	return err
}

func checkPrefix(s *discordgo.Session, e *discordgo.MessageCreate) (bool, int) {
	if strings.HasPrefix(e.Content, conf.Prefix) {
		return true, len(conf.Prefix)
	} else if strings.HasPrefix(e.Content, mentionPrefix) {
		return true, len(mentionPrefix)
	}
	return false, 0
}
