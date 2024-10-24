package service

import (
	"crypto/rand"
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
	GetRaw(key string) (data model.DBRecord, err error)
	Save(data model.Model) (err error)
	SaveRaw(data model.DBRecord) (err error)
	Delete(key string) (err error)
	GetToken() (token string, err error)
	ChangePasswd() (err error)
}

var _ Service = (*service)(nil)

type service struct {
	r *storage.Storage
}

func NewService(r *storage.Storage) *service {
	return &service{r: r}
}

func (s *service) ChangePasswd() (err error) {
	var token string
	isNew := cfg.User.GetString("packed_key") == ""
	bakToken := cfg.User.GetString("encryption_key")
	cfg.User.Set("encryption_key", nil)
	defer func() {
		if cfg.User.GetString("encryption_key") == "" && bakToken != "" {
			cfg.User.Set("encryption_key", bakToken)
		}
	}()
	token, err = s.GetToken()
	if err != nil {
		return
	}
	if isNew {
		return
	}
	var passRaw string
	passRaw, err = password.GetRawPass(true, cfg.PromptNewMasterPs, cfg.PromptConfirmMasterPs)
	userName := cfg.User.GetString("name")

	var (
		packedBytes  []byte
		cryptKeyPass = userName + string([]byte{9}) + passRaw
	)

	packedBytes, err = crypt.Encode([]byte(token), cryptKeyPass)
	if err != nil {
		return
	}
	cfg.User.Set("packed_key", hex.EncodeToString(packedBytes))

	return
}

func (s *service) GetToken() (token string, err error) {
	token = cfg.User.GetString("encryption_key")
	if token == "" {
		packed := cfg.User.GetString("packed_key")
		userName := cfg.User.GetString("name")
		var passRaw string
		if packed == "" {
			passRaw, err = password.GetRawPass(true, cfg.PromptNewMasterPs, cfg.PromptConfirmMasterPs)
		} else {
			passRaw, err = password.GetRawPass(false, cfg.PromptMasterPs, cfg.PromptConfirmMasterPs)
		}
		if err != nil {
			return
		}
		var (
			packedBytes, tokenBytes []byte
			cryptKeyPass            = userName + string([]byte{9}) + passRaw
		)
		if packed == "" {
			// fmt.Println("Creating new token... ")
			tokenBytes = make([]byte, 128)
			_, err = rand.Read(tokenBytes)
			if err != nil {
				err = errors.Join(errors.New("error create new token"), err)
				return
			}
			packedBytes, err = crypt.Encode(tokenBytes, cryptKeyPass)
			if err != nil {
				return
			}

			cfg.User.Set("packed_key", hex.EncodeToString(packedBytes))
		} else {
			packedBytes, err = hex.DecodeString(packed)
			if err != nil {
				err = errors.Join(errors.New("error hex.DecodeString"), err)
				return
			}
			tokenBytes, err = crypt.Decode(packedBytes, cryptKeyPass)
			if err != nil {
				err = errors.Join(errors.New("error decode token"), err)
				return
			}
		}
		token = string(tokenBytes)
		// cache token in config, it must be excluded from saving
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
	data.Items, err = s.r.DB.List(query)
	if err != nil {
		return
	}
	return
}

func (s *service) Get(key string) (data out.Item, err error) {
	var (
		r model.DBRecord
	)
	if r, err = s.GetRaw(key); err != nil {
		return
	}

	data.DBItem = r.DBItem
	var (
		deCrypted []byte
		token     string
	)
	token, err = s.GetToken()
	if err != nil {
		if !cfg.Glob.GetBool("debug") {
			err = errs.ErrPassword
		}
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

func (s *service) GetRaw(key string) (data model.DBRecord, err error) {
	if data, err = s.r.DB.Get(key); err != nil {
		return
	}
	if len(data.Blob) == 0 && data.Filename != nil {
		if data.Blob, err = s.r.File.GetStored(*data.Filename); err != nil {
			return
		}
	}
	return
}

func (s *service) Save(data model.Model) (err error) {
	if err = data.Validate(); err != nil {
		return
	}
	var r model.DBRecord
	r.Key = data.GetKey()
	r.Description = data.GetDescription()
	var blob []byte
	blob, err = model.NewPackedBytes(data)
	if err != nil {
		return
	}
	var token string
	token, err = s.GetToken()
	if err != nil {
		if !cfg.Glob.GetBool("debug") {
			err = errors.New("wrong password")
		}
		return
	}
	if r.Blob, err = crypt.Encode(blob, token); err != nil {
		return
	}
	err = s.SaveRaw(r)

	return
}

func (s *service) SaveRaw(data model.DBRecord) (err error) {
	if len(data.Blob) > cfg.MaxBlobSize {
		fileName := time.Now().Format("20060102150405") + "-" + data.Key
		err = s.r.File.SaveStore(fileName, data.Blob)
		if err != nil {
			return
		}
		data.Filename = &fileName
		data.Blob = nil
	}
	err = s.r.DB.Save(data)

	return
}

func (s *service) Delete(key string) (err error) {
	var r model.DBRecord
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
