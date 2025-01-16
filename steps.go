package gadget

import (
    "fmt"
    "os"
    "path/filepath"
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
)

// Step represents an action to be performed with its arguments.
type Step struct {
    Action Action
    Path   string
    Value  string
}

// Steps is a slice of Step.
type Steps []Step

// Run executes the step.
func (s Step) Run() error {
    switch s.Action {
    case Mkdir, MkdirCreateOnly:
        return os.MkdirAll(s.Path, 0775)
    case Rmdir:
        return os.Remove(s.Path)
    case Write, WriteBinary:
        if s.Value != "" {
            return os.WriteFile(s.Path, []byte(s.Value), 0664)
        }
        return nil
    case Remove:
        return os.Remove(s.Path)
    case Symlink:
        return os.Symlink(s.Path, s.Value)
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
        s.Value = filepath.Join(path, s.Value)
    }
    s.Path = filepath.Join(path, s.Path)
    return s
}

// Add appends a step to the Steps slice.
func (steps *Steps) Add(s Step) {
    *steps = append(*steps, s)
}

// AddSteps appends multiple steps to the Steps slice.
func (steps *Steps) AddSteps(more Steps) {
    *steps = append(*steps, more...)
}

// PrependPath prepends a path to all steps in the Steps slice.
func (steps Steps) PrependPath(path string) Steps {
    for i, s := range steps {
        steps[i] = s.PrependPath(path)
    }
    return steps
}

// Action represents the type of action to be performed in a step.
type Action int