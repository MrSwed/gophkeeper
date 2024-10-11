package model

import (
	pb "gophKeeper/internal/proto"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DBItem struct {
	Key         string     `db:"key" json:"key"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	SyncAt      *time.Time `db:"sync_at" json:"sync_at"`
}

type DBRecord struct {
	DBItem
	Filename *string `db:"filename,omitempty"`
	Blob     []byte  `db:"blob" json:"blob"`
}

func (d *DBRecord) FromItemSync(p *pb.ItemSync) {
	d.Key = p.Key
	d.Description = p.Description
	if p.UpdatedAt.IsValid() {
		d.UpdatedAt = new(time.Time)
		*d.UpdatedAt = p.UpdatedAt.AsTime()
	}
	if p.CreatedAt.IsValid() {
		d.CreatedAt = p.CreatedAt.AsTime()
	}
	d.Blob = p.Blob
	d.SyncAt = &[]time.Time{time.Now()}[0]
}

func (d *DBRecord) ToItemSync() (p *pb.ItemSync) {
	p = &pb.ItemSync{
		Key:         d.Key,
		Description: d.Description,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		Blob:        d.Blob,
	}
	if d.UpdatedAt != nil {
		p.UpdatedAt = timestamppb.New(*d.UpdatedAt)
	}
	return
}
