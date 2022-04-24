package impl

import (
	"context"
	"database/sql"

	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

const roomsTableName = "rooms"
const roomsTypeTableName = "rooms_type"
const roomsPriceTableName = "rooms_price"

type RoomTable struct {
	bun.BaseModel `bun:"table:rooms"` // table name rooms
	RoomNumber    uint64              `bun:"room_no,pk"`
	Description   string              `bun:"description"`
}

type RoomTypeTable struct {
	bun.BaseModel `bun:"table:rooms_type"` // table name rooms_type
	RoomType      string                   `bun:"model,unique:model_value"`
	RoomNumber    uint64                   `bun:"room_no,unique:model_value"` // FK from rooms table.
}

type RoomPriceTable struct {
	bun.BaseModel `bun:"table:rooms_price"` // table name rooms_price
	RoomPrice     float64                   `bun:"price,unique:price_value"`
	RoomNumber    uint64                    `bun:"room_no,unique:price_value"` // FK from rooms table.
}

func CreateInventoryTables(ctx context.Context, dbHandle *PostgresDB) error {
	db := dbHandle.DB(ctx)

	// first create the rooms table and then the rooms_type table.
	_, err := db.NewCreateTable().Model((*RoomTable)(nil)).
		IfNotExists().Exec(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": roomsTableName,
		}).WithError(err).Error("failed to create rooms table.")
		return err
	}

	_, err = db.NewCreateTable().Model((*RoomTypeTable)(nil)).
		IfNotExists().
		WithForeignKeys().ForeignKey(`("room_no") REFERENCES "rooms" ("room_no") ON DELETE CASCADE`).
		Exec(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": roomsTypeTableName,
		}).WithError(err).Error("failed to create rooms-type table.")
	}

	_, err = db.NewCreateTable().Model((*RoomPriceTable)(nil)).
		IfNotExists().
		WithForeignKeys().ForeignKey(`("room_no") REFERENCES "rooms" ("room_no") ON DELETE CASCADE`).
		Exec(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": roomsPriceTableName,
		}).WithError(err).Error("failed to create rooms-price table.")
	}
	return err
}

func DropInventoryTables(ctx context.Context, dbHandle *PostgresDB) error {
	db := dbHandle.DB(ctx)
	_, err := db.NewDropTable().Model((*RoomTypeTable)(nil)).IfExists().Exec(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": roomsTypeTableName,
		}).WithError(err).Error("failed to drop table rooms-type.")
		return err
	}
	_, err = db.NewDropTable().Model((*RoomTable)(nil)).IfExists().Exec(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": roomsTableName,
		}).WithError(err).Error("failed to drop table rooms-type.")
	}
	return err

}

func NewRoom(ctx context.Context, dbHandle *PostgresDB, roomNum uint64, roomDesc, roomType string) error {
	room := &RoomTable{
		RoomNumber:  roomNum,
		Description: roomDesc,
	}
	roomTypeRow := &RoomTypeTable{
		RoomType:   roomType,
		RoomNumber: roomNum,
	}

	db := dbHandle.DB(ctx)
	err := db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(room).Exec(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"db-table": roomsTableName,
				"rooom-no": roomNum,
				"Desc":     roomDesc,
			}).WithError(err).Error("failed to add new room to the table., cant execute transaction operation")
			return err
		}
		_, err = tx.NewInsert().Model(roomTypeRow).Exec(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"db-table":  roomsTypeTableName,
				"rooom-no":  roomNum,
				"Desc":      roomDesc,
				"room-type": roomType,
			}).WithError(err).Error("failed to add new room to the room-type table., cant execute transaction operation")
			return err
		}
		return nil
	})
	return err
}

func DeleteRoom(ctx context.Context, dbHandle *PostgresDB, roomNum uint64) error {
	db := dbHandle.DB(ctx)
	err := db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().Model((*RoomTable)(nil)).Where("room_no = ?", roomNum).Exec(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"db-table": roomsTableName,
				"room-no":  roomNum,
			}).WithError(err).Error("failed to delete the entry from rooms table!")
			return err
		}
		// no need to delete from types table as we created the FK to cascade the delete.
		return nil
	})
	return err
}
