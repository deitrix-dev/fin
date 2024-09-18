package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"slices"

	"github.com/deitrix/fin"
)

type Store struct {
	file string
}

type db struct {
	RecurringPayments []fin.RecurringPayment `json:"recurringPayments"`
}

func NewStore(file string) *Store {
	return &Store{file: file}
}

func (s *Store) RecurringPayment(_ context.Context, id string) (fin.RecurringPayment, error) {
	d, err := s.read()
	if err != nil {
		return fin.RecurringPayment{}, err
	}
	for _, rp := range d.RecurringPayments {
		if rp.ID == id {
			return rp, nil
		}
	}
	return fin.RecurringPayment{}, errors.New("recurring payment not found")
}

func (s *Store) RecurringPayments(_ context.Context) ([]fin.RecurringPayment, error) {
	d, err := s.read()
	if err != nil {
		return nil, err
	}
	return d.RecurringPayments, nil
}

func (s *Store) CreateRecurringPayment(_ context.Context, rp fin.RecurringPayment) error {
	d, err := s.read()
	if err != nil {
		return err
	}
	d.RecurringPayments = append(d.RecurringPayments, rp)
	return s.write(d)
}

func (s *Store) UpdateRecurringPayment(_ context.Context, rp fin.RecurringPayment) error {
	d, err := s.read()
	if err != nil {
		return err
	}
	index := slices.IndexFunc(d.RecurringPayments, func(rp2 fin.RecurringPayment) bool {
		return rp2.ID == rp.ID
	})
	if index == -1 {
		return errors.New("recurring payment not found")
	}
	rp.ID = d.RecurringPayments[index].ID
	d.RecurringPayments[index] = rp
	return s.write(d)
}

func (s *Store) DeleteRecurringPayment(_ context.Context, id string) error {
	d, err := s.read()
	if err != nil {
		return err
	}
	d.RecurringPayments = slices.DeleteFunc(d.RecurringPayments, func(rp fin.RecurringPayment) bool {
		return rp.ID == id
	})
	return s.write(d)
}

func (s *Store) read() (db, error) {
	f, err := os.Open(s.file)
	if errors.Is(err, os.ErrNotExist) {
		return db{}, nil
	}
	if err != nil {
		return db{}, err
	}
	defer f.Close()

	var d db
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return db{}, err
	}
	return d, nil
}

func (s *Store) write(d db) error {
	f, err := os.Create(s.file)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(d); err != nil {
		return err
	}
	return nil
}
