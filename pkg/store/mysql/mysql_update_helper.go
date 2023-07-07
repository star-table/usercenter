package mysql

import (
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (c Client) Update(obj Domain) error {
	conn, err := c.GetConnect()
	if err != nil {
		return err
	}
	err = conn.Collection(obj.TableName()).UpdateReturning(obj)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) TransUpdate(tx sqlbuilder.Tx, obj Domain) error {
	err := tx.Collection(obj.TableName()).UpdateReturning(obj)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) UpdateSmartWithCond(table string, cond db.Cond, upd Upd) (int64, error) {
	conn, err := c.GetConnect()
	if err != nil {
		return 0, err
	}
	res, err := conn.Update(table).Set(upd).Where(cond).Exec()
	if err != nil {
		Info(err)
		return 0, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return row, nil
}

func (c Client) UpdateSmart(table string, id int64, upd Upd) error {
	_, err := c.UpdateSmartWithCond(table, db.Cond{
		"id": id,
	}, upd)
	return err
}

func (c Client) TransUpdateSmartWithCond(tx sqlbuilder.Tx, table string, cond db.Cond, upd Upd) (int64, error) {
	res, err := tx.Update(table).Set(upd).Where(cond).Exec()
	if err != nil {
		Info(err)
		return 0, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		Info(err)
		return 0, err
	}

	return row, nil
}

func (c Client) TransUpdateSmart(tx sqlbuilder.Tx, table string, id int64, upd Upd) error {
	_, err := c.TransUpdateSmartWithCond(tx, table, db.Cond{
		"id": id,
	}, upd)
	return err
}
