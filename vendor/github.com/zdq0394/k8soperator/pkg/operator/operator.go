package operator

import (
	"fmt"
	"sync"

	"github.com/zdq0394/k8soperator/pkg/operator/controller"
)

// Operator is a controller, at code level have almost same contract of behavior
// but at a higher level it need to initialize some resource(usually CRDs) before
// start its execution.
type Operator interface {
	Initialize() error
	controller.Controller
}

type simpleOperator struct {
	crd         controller.CRD
	controller  controller.Controller
	initialized bool
	running     bool
	stateMu     sync.Mutex
}

// NewSimpleOperator create new instance of SimpleOperator.
func NewSimpleOperator(crd controller.CRD, ctrl controller.Controller) Operator {
	return &simpleOperator{
		crd:        crd,
		controller: ctrl,
	}
}

func (s *simpleOperator) Initialize() error {
	if s.isInitialized() {
		return nil
	}
	err := s.crd.Initialize()
	if err != nil {
		return err
	}
	s.setInitialized(true)
	return nil
}

func (s *simpleOperator) Run(stopC <-chan struct{}) error {
	if s.isRunning() {
		return fmt.Errorf("operator is already running")
	}
	s.setRunning(true)
	defer s.setRunning(false)
	if err := s.Initialize(); err != nil {
		return err
	}
	s.controller.Run(stopC)
	return nil
}

func (s *simpleOperator) isInitialized() bool {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	return s.initialized
}

func (s *simpleOperator) setInitialized(value bool) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.initialized = value
}

func (s *simpleOperator) isRunning() bool {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	return s.running
}

func (s *simpleOperator) setRunning(value bool) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.running = value
}
