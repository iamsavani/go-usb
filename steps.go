package gadget

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// Action types for different step actions.
const (
	Noop Action = iota
	Comment
	Mkdir
	MkdirCreateOnly
	Rmdir
	Write
	WriteBinary
	Remove
	Symlink
	Unlink
)

// Step represents an action to be performed with its arguments.
type Step struct {
	Action Action
	Arg0   string
	Arg1   string
}

// Steps is a slice of Step.
type Steps []Step

// Action represents the type of action to be performed in a step.
type Action int

// Run executes the step.
func (s Step) Run() error {
	switch s.Action {
	case Mkdir, MkdirCreateOnly:
		return os.MkdirAll(s.Arg0, 0775)
	case Rmdir:
		return os.Remove(s.Arg0)
	case Write, WriteBinary:
		if s.Arg1 != "" {
			return os.WriteFile(s.Arg0, []byte(s.Arg1), 0664)
		}
		return nil
	case Remove:
		return os.Remove(s.Arg0)
	case Symlink:
		return os.Symlink(s.Arg0, s.Arg1)
	default:
		return nil
	}
}

// Run executes all steps in the Steps slice.
func (steps Steps) Run() error {
	for i, s := range steps {
		if err := s.Run(); err != nil {
			return fmt.Errorf("step %d %+v: error %w", i, s, err)
		}
	}
	return nil
}

// PrependPath prepends a path to the step's arguments.
func (s Step) PrependPath(path string) Step {
	switch s.Action {
	case Noop, Comment:
		return s
	case Symlink:
		s.Arg1 = filepath.Join(path, s.Arg1)
	}
	s.Arg0 = filepath.Join(path, s.Arg0)
	return s
}

// Add appends a step to the Steps slice.
func (ss *Steps) Append(s Step) Steps {
	*ss = append(*ss, s)
	return *ss
}

// AddSteps appends multiple steps to the Steps slice.
func (ss *Steps) Extend(more Steps) Steps {
	*ss = append(*ss, more...)
	return *ss
}

// PrependPath prepends a path to all steps in the Steps slice.
func (steps Steps) PrependPath(path string) Steps {
	for i, s := range steps {
		steps[i] = s.PrependPath(path)
	}
	return steps
}

// undo generates the step to undo the current step.
func (s Step) Undo() Step {
	switch s.Action {
	case Mkdir:
		return Step{Rmdir, s.Arg0, ""}
	case Symlink:
		return Step{Remove, s.Arg1, ""}
	default:
		return Step{Noop, "", ""}
	}
}

// Clone returns a copy of the steps.
func (steps Steps) Clone() Steps { return slices.Clone(steps) }

// Undo generates the steps to undo the current steps.
func (steps Steps) Reverse() (rev Steps) {
	rev = steps.Clone()
	slices.Reverse(rev)
	return
}

// Reverse reverses the order of the steps.
func (steps Steps) Undo() Steps {
	var undo = steps.Clone()
	for i := range undo {
		undo[i] = undo[i].Undo()
	}
	return undo
}
