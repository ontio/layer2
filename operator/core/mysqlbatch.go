package core

import (
	"fmt"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlInsertBatch struct {
	db                 *sql.DB
	valueStrings       []string
	valueArgs          []interface{}
	valueStr           string
	stmt               string
	cols               int
	counter            int
}

func NewMysqlInsertBatch(db *sql.DB, cols int, valueStr string, stmt string) *MysqlInsertBatch {
	batch := &MysqlInsertBatch{
		db: db,
		valueStrings: make([]string, 0, 10000),
		valueArgs: make([]interface{}, 0, cols * 10000),
		valueStr: valueStr,
		stmt: stmt,
		cols: cols,
		counter: 0,
	}
	return batch
}

func (this *MysqlInsertBatch) Insert(data []interface{}) error {
	if len(data) != this.cols {
		return fmt.Errorf("insert data cols is not right!")
	}
	this.valueStrings = append(this.valueStrings, this.valueStr)
	this.valueArgs = append(this.valueArgs, data...)
	this.counter ++

	if this.counter >= 10000 {
		err := this.Commit()
		if err != nil {
			return err
		}
		this.valueStrings = this.valueStrings[0:0]
		this.valueArgs = this.valueArgs[0:0]
		this.counter = 0
	}
	return nil
}

func (this *MysqlInsertBatch) Commit() error {
	stmt := fmt.Sprintf("%s VALUES %s", this.stmt, strings.Join(this.valueStrings, ","))
	_, dberr := this.db.Exec(stmt, this.valueArgs...)
	if dberr != nil {
		return fmt.Errorf("batch commit error: %s", dberr.Error())
	} else {
		return nil
	}
}

func (this *MysqlInsertBatch) Close() error {
	err := this.Commit()
	if err != nil {
		return err
	}
	this.valueStrings = this.valueStrings[0:0]
	this.valueArgs = this.valueArgs[0:0]
	this.counter = 0
	return nil
}

type MysqlUpdateBatch struct {
	db                 *sql.DB
	valueStrings       []string
	valueArgs          []interface{}
	valueStr           string
	stmt               string
	update             string
	cols               int
	counter            int
}

func NewMysqlUpdateBatch(db *sql.DB, cols int, valueStr string, stmt string, update string) *MysqlUpdateBatch {
	batch := &MysqlUpdateBatch{
		db: db,
		valueStrings: make([]string, 0, 10000),
		valueArgs: make([]interface{}, 0, cols * 10000),
		valueStr: valueStr,
		stmt: stmt,
		update: update,
		cols: cols,
		counter: 0,
	}
	return batch
}

func (this *MysqlUpdateBatch) Insert(data []interface{}) error {
	if len(data) != this.cols {
		return fmt.Errorf("insert data cols is not right!")
	}
	this.valueStrings = append(this.valueStrings, this.valueStr)
	this.valueArgs = append(this.valueArgs, data...)
	this.counter ++

	if this.counter >= 10000 {
		err := this.Commit()
		if err != nil {
			return err
		}
		this.valueStrings = this.valueStrings[0:0]
		this.valueArgs = this.valueArgs[0:0]
		this.counter = 0
	}
	return nil
}

func (this *MysqlUpdateBatch) Commit() error {
	stmt := fmt.Sprintf("%s VALUES %s %s", this.stmt, strings.Join(this.valueStrings, ","), this.update)
	_, dberr := this.db.Exec(stmt, this.valueArgs...)
	if dberr != nil {
		return fmt.Errorf("batch commit error: %s", dberr.Error())
	} else {
		return nil
	}
}

func (this *MysqlUpdateBatch) Close() error {
	err := this.Commit()
	if err != nil {
		return err
	}
	this.valueStrings = this.valueStrings[0:0]
	this.valueArgs = this.valueArgs[0:0]
	this.counter = 0
	return nil
}
