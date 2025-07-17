package db

import (
	"Go_Team00.ID_376234-Team_TL_barievel/internal/entities"
	"time"
)

type EntryDB struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	SessionID string    `gorm:"index;not null"`
	Frequency float64   `gorm:"not null"`
	Timestamp time.Time `gorm:"type:timestamptz;not null"`
}

func EntryToEntryDB(entry entities.Entry) *EntryDB {
	return &EntryDB{
		ID:        0,
		SessionID: entry.SessionId,
		Frequency: entry.Frequency,
		Timestamp: entry.Timestamp,
	}
}
func EntryDbToEntry(entryDB EntryDB) entities.Entry {
	return entities.Entry{
		SessionId: entryDB.SessionID,
		Frequency: entryDB.Frequency,
		Timestamp: entryDB.Timestamp,
	}
}
