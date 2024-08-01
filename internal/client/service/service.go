package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"gophKeeper/internal/client/input/password"
	"time"

	cfg "gophKeeper/internal/client/config"
	"gophKeeper/internal/client/crypt"
	errs "gophKeeper/internal/client/errors"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/out"
	"gophKeeper/internal/client/storage"
)

type Service interface {
	List(query model.ListQuery) (data out.List, err error)
	Get(key string) (data out.Item, err error)
	Save(data model.Model) (err error)
	Delete(key string) (err error)
	GetToken() (token string, err error)
}

var _ Service = (*service)(nil)

type service struct {
	r *storage.Storage
}

func NewService(r *storage.Storage) *service {
	return &service{r: r}
}

func (s *service) GetToken() (token string, err error) {
	token = cfg.User.GetString("encryption_key")
	if token == "" {
		packed := cfg.User.GetString("packed_key")
		userName := cfg.User.GetString("name")
		var passRaw string
		passRaw, err = password.GetRawPass(packed == "")
		if err != nil {
			return
		}
		var (
			packedBytes, tokenBytes []byte
		)
		if packed == "" {
			// fmt.Println("Creating new token... ")
			tokenBytes = make([]byte, 128)
			_, err = rand.Read(tokenBytes)
			if err != nil {
				err = errors.Join(errors.New("error create new token"), err)
				return
			}
			tokenStr := hex.EncodeToString(tokenBytes)

			packedBytes, err = crypt.Encode([]byte(tokenStr), userName+passRaw)
			if err != nil {
				return
			}
			cfg.User.Set("packed_key", packedBytes)
		} else {
			tokenBytes, err = crypt.Decode([]byte(packed), userName+passRaw)
			if err != nil {
				err = errors.Join(errors.New("error decode token"), err)
				return
			}
		}
		token = string(tokenBytes)
		// cache token in config, ot must be excluded from saving
		cfg.User.Set("encryption_key", token)
	}
	return
}

func (s *service) List(query model.ListQuery) (data out.List, err error) {
	if err = query.Validate(); err != nil {
		return
	}
	if data.Total, err = s.r.DB.Count(query); err != nil {
		return
	}
	var items []storage.DBItem
	items, err = s.r.DB.List(query)
	if err != nil {
		return
	}
	data.FromDBItems(items...)
	return
}

func (s *service) Get(key string) (data out.Item, err error) {
	var (
		r storage.DBRecord
	)
	if r, err = s.r.DB.Get(key); err != nil {
		return
	}
	if len(r.Blob) == 0 && r.Filename != nil {
		if r.Blob, err = s.r.File.GetStored(*r.Filename); err != nil {
			return
		}
	}
	data.DBItem = r.DBItem
	var (
		deCrypted []byte
		token     string
	)
	token, err = s.GetToken()
	if err != nil {
		return
	}
	deCrypted, err = crypt.Decode(r.Blob, token)
	if err != nil {
		err = errs.ErrDecode
		return
	}
	err = json.Unmarshal(deCrypted, &data)
	if dataSan, ok := data.Data.(model.Sanitisable); ok {
		dataSan.Sanitize()
	}
	return
}

func (s *service) Save(data model.Model) (err error) {
	if err = data.Validate(); err != nil {
		return
	}
	var r storage.DBRecord
	if r, err = s.r.DB.Get(data.GetKey()); err != nil &&
		!errors.Is(err, sql.ErrNoRows) {
		return
	}
	if r.Key == "" {
		r.Key = data.GetKey()
	}
	r.Description = data.GetDescription()

	var blob []byte
	blob, err = data.Bytes()
	if err != nil {
		return
	}
	var token string
	token, err = s.GetToken()
	if err != nil {
		return
	}

	if blob, err = crypt.Encode(blob, token); err != nil {
		return
	}

	if len(blob) > cfg.MaxBlobSize {
		fileName := time.Now().Format("20060102150405") + "-" + r.Key
		err = s.r.File.SaveStore(fileName, blob)
		if err != nil {
			return
		}
		r.Filename = &fileName
	} else {
		r.Blob = blob
	}

	err = s.r.DB.Save(r)

	return
}

func (s *service) Delete(key string) (err error) {
	var r storage.DBRecord
	if r, err = s.r.DB.Get(key); err != nil {
		return
	}
	if r.Filename != nil && *r.Filename != "" {
		if err = s.r.File.Delete(*r.Filename); err != nil {
			return
		}
	}
	err = s.r.DB.Delete(key)
	return
}
