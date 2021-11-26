package apbun

import (
	"errors"
	"testing"
	"time"
)

type Step1 struct{}

func (Step1) EXQ() error  { println("step1 exq"); return nil }
func (Step1) Undo() error { println("step1 undo"); return nil }

type Step2 struct{}

func (Step2) EXQ() error  { println("step2 exq"); return nil }
func (Step2) Undo() error { println("step2 undo"); return nil }

type Step3 struct{}

func (Step3) EXQ() error  { return errors.New("step3 error") }
func (Step3) Undo() error { println("step3 undo"); return nil }

type Step4 struct{}

func (Step4) EXQ() error  { println("step1 exqed"); return nil }
func (Step4) Undo() error { println("step4 undo"); return nil }

type StepTimeOut struct{}

func (StepTimeOut) EXQ() error  { time.Sleep(2 * time.Second); println("step1 exqed"); return nil }
func (StepTimeOut) Undo() error { println("step_time_out undo"); return nil }

func TestAp(t *testing.T) {
	ap := NewBun(0)
	ap.AppendCommands(Step1{}, Step2{}, Step3{}, Step4{})
	if err := ap.AP(); err == nil || err.Error() != "step3 error" {
		t.Errorf("err should not be nil, but it is %s", err.Error())
	}
}
func TestApWithTimeOut(t *testing.T) {
	ap := NewBun(1)
	ap.AppendCommands(Step1{}, Step2{}, StepTimeOut{}, Step4{})
	if err := ap.AP(); err == nil || err.Error() != "ap_bun time out" {
		t.Errorf("err should not be nil, but it is %s", err.Error())
	}
}
