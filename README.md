# APBUN

APBUN is a Go package that supports a easy way to warp and undo commands.

By the way, APBUN means 'Pineapple bun' that I used to love, so it should be 'undo' now.

Getting Started
===============

## Installing

To start using APBUN, install Go and run `go get`:

```sh
$ go get -u github.com/projectred/apbun
```

## Run

```golang
type Step1 struct{}

func (Step1) EXQ() error  { println("step1 exq"); return nil }
func (Step1) Undo() error { println("step1 undo"); return nil }

type Step2 struct{}

func (Step2) EXQ() error  { println("step2 exq"); return nil }
func (Step2) Undo() error { println("step2 undo"); return nil }

type StepWithError struct{}

func (StepWithError) EXQ() error  { return errors.New("step3 error") }
func (StepWithError) Undo() error { println("step3 undo"); return nil }

type Step4 struct{ NoUndo }

func (Step4) EXQ() error { println("step1 exqed"); return nil }

func TestAp(t *testing.T) {
	ap := NewBun(0)
	ap.AppendCommands(Step1{}, Step2{}, StepWithError{}, Step4{})
	if err := ap.AP(); err == nil || err.Error() != "step3 error" {
		t.Errorf("err should not be nil, but it is %s", err.Error())
	}
}

// Output:
// === RUN   TestAp
// step1 exq
// step2 exq
// step2 undo
// step1 undo



type StepTimeOut struct{}

func (StepTimeOut) EXQ() error  { time.Sleep(2 * time.Second); println("step exqed"); return nil }
func (StepTimeOut) Undo() error { println("step_time_out undo"); return nil }

func TestApWithTimeOut(t *testing.T) {
	ap := NewBun(1)
	ap.AppendCommands(Step1{}, Step2{}, StepTimeOut{}, Step4{})
	if err := ap.AP(); err == nil || err.Error() != "ap_bun time out" {
		t.Errorf("err should not be nil, but it is %s", err.Error())
	}
}

// output: 
// === RUN   TestApWithTimeOut
// step1 exq
// step2 exq
// step2 undo
// step1 undo
```