package mysql

import (
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (c Client) SelectById(table string, id interface{}, obj interface{}) error {
	conn, err := c.GetConnect()
	if err != nil {
		return err
	}
	err = conn.Collection(table).Find(db.Cond{
		"id":        id,
		"is_delete": 2,
	}).One(obj)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) SelectCountByCond(table string, cond db.Cond, unions ...*db.Union) (uint64, error) {
	conn, err := c.GetConnect()
	if err != nil {
		return 0, err
	}
	mid := conn.Collection(table).Find(cond)
	if len(unions) > 0 {
		for _, union := range unions {
			mid = mid.And(union)
		}
	}
	unit, err := mid.Count()
	if err != nil {
		return 0, err
	}
	return unit, nil
}

func (c Client) TransSelectCountByCond(tx sqlbuilder.Tx, table string, cond db.Cond) (uint64, error) {
	unit, err := tx.Collection(table).Find(cond).Count()
	if err != nil {
		return 0, err
	}
	return unit, nil
}

func (c Client) SelectOneByCond(table string, cond db.Cond, obj interface{}) error {
	conn, err := c.GetConnect()
	if err != nil {
		return err
	}
	err = conn.Collection(table).Find(cond).One(obj)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) TransSelectOneByCond(tx sqlbuilder.Tx, table string, cond db.Cond, obj interface{}) error {
	err := tx.Collection(table).Find(cond).One(obj)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) SelectByQuery(query string, objs interface{}, args ...interface{}) error {
	conn, err := c.GetConnect()
	if err != nil {
		return err
	}
	var iter sqlbuilder.Iterator = nil
	if len(args) > 0 {
		iter = conn.Iterator(query, args...)
	} else {
		iter = conn.Iterator(query)
	}
	err = iter.All(objs)
	return err
}

func (c Client) TransSelectByQuery(tx sqlbuilder.Tx, query string, objs interface{}, args ...interface{}) error {
	var iter sqlbuilder.Iterator = nil
	if len(args) > 0 {
		iter = tx.Iterator(query, args...)
	} else {
		iter = tx.Iterator(query)
	}
	err := iter.All(objs)
	return err
}

func (c Client) SelectAllByCond(table string, cond db.Cond, objs interface{}) error {
	conn, err := c.GetConnect()
	if err != nil {
		return err
	}
	err = conn.Collection(table).Find(cond).All(objs)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) SelectAllByCondWithColumns(table string, columns interface{},cond db.Cond, objs interface{}) error {
	conn, err := c.GetConnect()
	if err != nil {
		return err
	}
	if columns == nil  {
		columns = db.Raw("*")
	}
	err = conn.Select(columns).From(table).Where(cond).All(objs)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) TransSelectAllByCond(tx sqlbuilder.Tx, table string, cond db.Cond, objs interface{}) error {
	err := tx.Collection(table).Find(cond).All(objs)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) SelectAllByCondWithPageAndOrder(table string, cond db.Cond, union *db.Union, page int, size int, order interface{}, objs interface{}) (uint64, error) {
	conn, err := c.GetConnect()
	if err != nil {
		return 0, err
	}

	mid := conn.Collection(table).Find(cond)
	if union != nil {
		mid = mid.And(union)
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	count, err := mid.TotalEntries()
	if err != nil {
		return 0, err
	}
	err = mid.All(objs)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c Client) SelectAllByCondWithPageAndOrderUnion(table string, cond db.Cond, page int, size int, order interface{}, objs interface{}, unions ...*db.Union) (uint64, error) {
	conn, err := c.GetConnect()
	if err != nil {
		return 0, err
	}
	conn.SetLogging(true) // for debug
	mid := conn.Collection(table).Find(cond)
	if len(unions) > 0 {
		for _, union := range unions {
			mid = mid.And(union)
		}
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	count, err := mid.TotalEntries()
	if err != nil {
		return 0, err
	}
	err = mid.All(objs)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func TransSelectAllByCondWithPageAndOrder(tx sqlbuilder.Tx, table string, cond db.Cond, union *db.Union, page int, size int, order interface{}, objs interface{}) (uint64, error) {
	mid := tx.Collection(table).Find(cond)
	if union != nil {
		mid = mid.And(union)
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	count, err := mid.TotalEntries()
	if err != nil {
		return 0, err
	}
	err = mid.All(objs)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c Client) SelectAllByCondWithNumAndOrder(table string, cond db.Cond, union *db.Union, page int, size int, order interface{}, objs interface{}) error {
	conn, err := c.GetConnect()
	if err != nil {
		return err
	}

	mid := conn.Collection(table).Find(cond)
	if union != nil {
		mid = mid.And(union)
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	err = mid.All(objs)
	if err != nil {
		return err
	}
	return nil
}

func TransSelectAllByCondWithNumAndOrder(tx sqlbuilder.Tx, table string, cond db.Cond, union *db.Union, page int, size int, order interface{}, objs interface{}) error {
	mid := tx.Collection(table).Find(cond)
	if union != nil {
		mid = mid.And(union)
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	err := mid.All(objs)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) IsExistByCond(table string, cond db.Cond) (bool, error) {
	conn, err := c.GetConnect()
	if err != nil {
		return false, err
	}
	exist, err := conn.Collection(table).Find(cond).Exists()
	if err != nil {
		return false, err
	}
	return exist, nil
}
