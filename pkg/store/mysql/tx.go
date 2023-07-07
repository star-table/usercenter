package mysql

import (
	"upper.io/db.v3/lib/sqlbuilder"
)

func (c Client) TransX(txFunc func(tx sqlbuilder.Tx) error) error {
	conn, err := c.GetConnect()
	if err != nil {
		Info(err)
		return err
	}
	tx, err := conn.NewTx(nil)
	if err != nil {
		Info(err)
		return err
	}
	defer CloseTx(conn, tx)

	err = txFunc(tx)
	if err != nil {
		Info(err)
		Rollback(tx)
		return err
	}

	err = tx.Commit()

	if err != nil {
		Info("tx.Commit(): ", err)
		return err
	}
	return nil
}
