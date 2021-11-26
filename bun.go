package apbun

import (
	"context"
	"errors"
	"log"
	"time"
)

type Bun struct {
	expireTime int64
	undons     []Undo
	commands   []Command
	log.Logger
}

func NewBun(expireTime int64) *Bun {
	bun := &Bun{expireTime: expireTime}
	return bun
}
func (b *Bun) AppendCommands(cs ...Command) { b.commands = append(b.commands, cs...) }

func (b *Bun) AP() error {
	ctx, timeOutF := b.withTimeout()
	defer timeOutF()
	for i := 0; i < len(b.commands); i++ {
		errc := GoWithError(func() error {
			if err := b.commands[i].EXQ(); err != nil {
				return err
			}
			b.undons = append(b.undons, b.commands[i])
			return nil
		})
		select {
		case <-ctx.Done():
			b.Undo()
			return errors.New("ap_bun time out")
		case err := <-errc:
			if err != nil {
				b.Undo()
				return err
			}
		}
	}
	return nil
}
func (b *Bun) Undo() error {
	ctx, timeOutF := b.withTimeout()
	defer timeOutF()
	for j := len(b.undons) - 1; j >= 0; j-- {
		errc := GoWithError(func() error { return b.undons[j].Undo() })
		select {
		case <-ctx.Done():
			// print error
			return nil
		case err := <-errc:
			if err != nil {
				// print error
			}
		}
	}
	return nil
}
func (b *Bun) withTimeout() (context.Context, func()) {
	if b.expireTime > 0 {
		return context.WithTimeout(context.Background(), time.Duration(b.expireTime)*time.Second)
	}
	return context.Background(), func() {}
}

func GoWithError(f func() error) chan error {
	errc := make(chan error)
	go func() {
		errc <- f()
	}()
	return errc
}

type Command interface {
	EXQ() error
	Undo
}

type Undo interface {
	Undo() error
}
