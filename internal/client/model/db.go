package model

import (
	"time"

	pb "gophKeeper/internal/proto"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DBItem struct {
	Key         string     `db:"key" json:"key"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	SyncAt      *time.Time `db:"sync_at" json:"sync_at"`
	Description string     `db:"description" json:"description"`
}

type DBRecord struct {
	DBItem
	Blob     []byte  `db:"blob" json:"blob"`
	Filename *string `db:"filename,omitempty"`
}

// IsDeleted checks if the DBRecord is considered deleted.
// A record is considered deleted if its Blob is empty and its Filename is nil.
func (d *DBRecord) IsDeleted() bool {
	return len(d.Blob) == 0 && d.Filename == nil
}

// FromItemSync
//
//	convert remote proto item to local db record
func (d *DBRecord) FromItemSync(p *pb.ItemSync) {
	d.Key = p.Key
	d.Description = p.Description
	if p.UpdatedAt.IsValid() {
		d.UpdatedAt = new(time.Time)
		*d.UpdatedAt = p.UpdatedAt.AsTime().Local()
	}
	if p.CreatedAt.IsValid() {
		d.CreatedAt = p.CreatedAt.AsTime().Local()
	}
	d.Blob = p.Blob
	d.SyncAt = &[]time.Time{time.Now()}[0]
}

// ToItemSync
//
//	Convert local db record to proto item for send to remote
//	we save to local sqlite datetime in localtime zone without zone ext "+03:00"
//	so we need to convert it to utc
func (d *DBRecord) ToItemSync() (p *pb.ItemSync) {
	_, z := time.Now().Zone()
	p = &pb.ItemSync{
		Key:         d.Key,
		Description: d.Description,
		CreatedAt:   timestamppb.New(d.CreatedAt.Add(-time.Duration(z) * time.Second)),
		Blob:        d.Blob,
	}
	if d.UpdatedAt != nil {
		p.UpdatedAt = timestamppb.New(d.UpdatedAt.Add(-time.Duration(z) * time.Second))
	}
	return
}
