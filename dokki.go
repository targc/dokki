package dokki

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"testing"
	"unsafe"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

var lock sync.Mutex

type Setup struct {
	comps []TestSuiteDockerCompose
}

func NewSetup() *Setup {
	return &Setup{}
}

type TestSuiteDockerCompose struct {
	composeStack compose.ComposeStack
	services     []string
}

func (s *Setup) Down(t *testing.T) {
	for i := range s.comps {
		c := s.comps[len(s.comps)-1-i]
		c.teardown(t)
	}
}

func (s *Setup) SetupTestSuiteDockerCompose(t *testing.T, composeFile string, services []string) *TestSuiteDockerCompose {
	composeStack, err := compose.NewDockerComposeWith(
		compose.WithStackFiles(composeFile),
		compose.WithLogger(log.Default()),
	)

	if err != nil {
		t.Fatalf("Failed to create compose stack: %v", err)
	}

	unsafe_setUnexportedField(composeStack, "name", "dokki")

	ts := TestSuiteDockerCompose{
		composeStack: composeStack,
		services:     services,
	}

	ts.up(t)

	s.comps = append(s.comps, ts)

	return &ts
}

func (ts *TestSuiteDockerCompose) up(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()

	err := ts.composeStack.Up(
		ctx,
		compose.Wait(true),
		compose.RunServices(
			ts.services...,
		),
	)

	if err != nil {
		t.Fatalf("Failed to start compose stack: %v", err)
	}
}

func (ts *TestSuiteDockerCompose) teardown(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	if ts.composeStack == nil {
		return
	}

	ctx := context.Background()

	err := ts.composeStack.Down(
		ctx,
		compose.RemoveOrphans(true),
		compose.RemoveVolumes(true),
		compose.RemoveImages(compose.RemoveImagesLocal),
	)

	if err != nil {
		t.Fatalf("Failed to down compose stack: %v", err)
	}
}

func unsafe_setUnexportedField(ptrToStruct any, field string, value any) error {

	v := reflect.ValueOf(ptrToStruct)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("expect pointer to struct")
	}

	s := v.Elem()

	f := s.FieldByName(field)

	if !f.IsValid() {
		return fmt.Errorf("no such field %q", field)
	}

	if !f.CanAddr() {
		return fmt.Errorf("field %q not addressable", field)
	}

	fp := unsafe.Pointer(f.UnsafeAddr())
	w := reflect.NewAt(f.Type(), fp).Elem()

	val := reflect.ValueOf(value)

	if !val.Type().AssignableTo(f.Type()) {
		if val.Type().ConvertibleTo(f.Type()) {
			val = val.Convert(f.Type())
		} else {
			return fmt.Errorf("cannot assign %s to %s", val.Type(), f.Type())
		}
	}

	w.Set(val)

	return nil
}
