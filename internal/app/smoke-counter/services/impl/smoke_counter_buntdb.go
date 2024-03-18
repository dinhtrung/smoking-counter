package impl

import (
	"errors"
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/services"
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/web/dto"
	"github.com/tidwall/buntdb"
	"log/slog"
	"strings"
	"time"
)

type SmokeServiceBuntDB struct {
	db *buntdb.DB
}

func NewSmokeServiceBuntDB(db *buntdb.DB) services.SmokeService {
	return &SmokeServiceBuntDB{
		db: db,
	}
}

func (s *SmokeServiceBuntDB) GetAll(user string) []*dto.DailySmokingDTO {
	var res []*dto.DailySmokingDTO
	if err := s.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(user+":*", func(key, value string) bool {
			events := strings.Split(value, ",")
			row := &dto.DailySmokingDTO{
				Date:   key[len(user)+1:],
				Events: events,
				Count:  len(events),
			}
			res = append(res, row)
			return true
		})
	}); err != nil {
		slog.Error("unable to retrieve data", "error", err, "user", user)
	}
	return res
}

func (s *SmokeServiceBuntDB) Create(user string, hour string) (*dto.DailySmokingDTO, error) {
	t, err := time.Parse("15:04", hour)
	if err != nil {
		return nil, err
	}
	dt := time.Now().Format("2006-01-02")
	key := user + ":" + dt
	val := t.Format("15:04")
	var events []string
	if err := s.db.Update(func(tx *buntdb.Tx) error {
		existsValue, err := tx.Get(key)
		if err != nil && !errors.Is(err, buntdb.ErrNotFound) {
			return err
		}
		events = strings.Split(existsValue, ",")
		events = append(events, val)
		val = strings.Join(events, ",")
		_, _, err = tx.Set(key, val, nil)
		return err
	}); err != nil {
		return nil, err
	}
	return &dto.DailySmokingDTO{
		Date:   dt,
		Events: events,
		Count:  len(events),
	}, nil
}

func (s *SmokeServiceBuntDB) Delete(user, hour string) error {
	t, err := time.Parse("15:04", hour)
	if err != nil {
		return err
	}
	dt := time.Now().Format("2006-01-02")
	key := user + ":" + dt
	var events []string
	if err := s.db.View(func(tx *buntdb.Tx) error {
		existsValue, err := tx.Get(key)
		if err != nil && !errors.Is(err, buntdb.ErrNotFound) {
			return err
		}
		events = strings.Split(existsValue, ",")
		return nil
	}); err != nil {
		return err
	}
	val := t.Format("15:04")
	var newEvents []string
	for _, event := range events {
		if event != val {
			newEvents = append(newEvents, event)
		}
	}
	if err := s.db.Update(func(tx *buntdb.Tx) error {
		val = strings.Join(newEvents, ",")
		_, _, err = tx.Set(key, val, nil)
		return err
	}); err != nil {
		return err
	}
	return nil
}

func (s *SmokeServiceBuntDB) DeleteAll(user string) error {
	var keys []string
	if err := s.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(user+":", func(key, value string) bool {
			keys = append(keys, key)
			return true
		})
	}); err != nil {
		return err
	}
	// delete all keys
	return s.db.Update(func(tx *buntdb.Tx) error {
		for _, key := range keys {
			if _, err := tx.Delete(key); err != nil && !errors.Is(err, buntdb.ErrNotFound) {
				return err
			}
		}
		return nil
	})
}

func (s *SmokeServiceBuntDB) DeleteOne(user, date, hour string) error {
	//TODO implement me
	panic("implement me")
}

func (s *SmokeServiceBuntDB) DeleteAllByDate(user, date string) error {
	//TODO implement me
	panic("implement me")
}

func (s *SmokeServiceBuntDB) DeleteAllBefore(user, date string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SmokeServiceBuntDB) DeleteAllAfter(user, date string) (int, error) {
	//TODO implement me
	panic("implement me")
}
