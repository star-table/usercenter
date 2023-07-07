package mysql

import (
	"upper.io/db.v3/lib/sqlbuilder"
)

func (c Client) Insert(obj Domain) error {
	conn, err := c.GetConnect()
	if err != nil {
		return err
	}
	_, err = conn.Collection(obj.TableName()).Insert(obj)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) InsertReturnId(obj Domain) (interface{}, error) {
	conn, err := c.GetConnect()
	if err != nil {
		return nil, err
	}
	id, err := conn.Collection(obj.TableName()).Insert(obj)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (c Client) TransInsert(tx sqlbuilder.Tx, obj Domain) error {
	_, err := tx.Collection(obj.TableName()).Insert(obj)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) TransInsertReturnId(tx sqlbuilder.Tx, obj Domain) (interface{}, error) {
	id, err := tx.Collection(obj.TableName()).Insert(obj)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (c Client) TransBatchInsert(tx sqlbuilder.Tx, obj Domain, objs []interface{}) error {

	//a := objs.([]interface{})

	batch := tx.InsertInto(obj.TableName()).Batch(len(objs))
	go func() {
		defer batch.Done()
		for i := range objs {
			batch.Values(objs[i])
		}
	}()
	err := batch.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (c Client) BatchInsert(obj Domain, objs []interface{}) error {
	conn, err := c.GetConnect()
	if err != nil {
		Info(err)
		return err
	}

	batch := conn.InsertInto(obj.TableName()).Batch(len(objs))
	go func() {
		defer batch.Done()
		for i := range objs {
			batch.Values(objs[i])
		}
	}()
	err = batch.Wait()
	if err != nil {
		return err
	}

	return nil
}

//func BatchDone(pos []interface{}, batch *sqlbuilder.BatchInserter) {
//		defer batch.Done()
//		for i := range pos {
//			batch.Values(pos[i])
//		}
//}
