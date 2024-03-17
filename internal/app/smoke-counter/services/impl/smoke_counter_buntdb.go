package impl

import (
	"errors"
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/services"
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/web/dto"
	"github.com/tidwall/buntdb"
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
	s.db.View(func(tx *buntdb.Tx) error {
		tx.AscendKeys(user+":", func(key, value string) bool {
			events := strings.Split(value, ",")
			row := &dto.DailySmokingDTO{
				Date:   key[len(user)+1:],
				Events: events,
				Count:  len(events),
			}
			res = append(res, row)
			return true
		})
		return nil
	})
	return res
}

func (s *SmokeServiceBuntDB) Create(user string, hour string) (*dto.DailySmokingDTO, error) {
	t, err := time.Parse("15:04", hour)
	if err != nil {
		return nil, err
	}
	key := user + ":" + time.Now().Format("2006-01-02")
	val := t.Format("15:04")
	s.db.Update(func(tx *buntdb.Tx) error {
		existsValue, err := tx.Get(key)
		if err != nil && !errors.Is(err, buntdb.ErrNotFound) {
			return err
		}
		events := strings.Split(existsValue, ",")
		events = append(events, val)
		_, _, err = tx.Set(key, strings.Join(events, ","), nil)
		return err
	})
	//TODO implement me
	panic("implement me")
}

func (s *SmokeServiceBuntDB) DeleteAll(user string) error {
	var keys []string
	if err := s.db.View(func(tx *buntdb.Tx) error {
		tx.AscendKeys(user+":", func(key, value string) bool {
			keys = append(keys, key)
			return true
		})
		return nil
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
