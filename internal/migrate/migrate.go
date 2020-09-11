package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Migrator struct
type Migrator struct {
	db        *sql.DB
	path      string
	tableName string
	upName    string
	downName  string
	ts        string
	tp        int
}

var (
	createType = 1
	alterType  = 2
)

// New return migrator instance
func New(db *sql.DB, path string) *Migrator {
	migrator := &Migrator{
		db:   db,
		path: path,
		ts:   time.Now().Format("20060102150405"),
	}
	if err := migrator.createMigrationTable(); err != nil {
		panic(err)
	}
	return migrator
}

// SetName set migration name
func (m *Migrator) SetName(name string) error {
	m.tableName = name
	m.upName = fmt.Sprintf("%s_%s.up.sql", m.ts, name)
	m.downName = fmt.Sprintf("%s_%s.down.sql", m.ts, name)
	s := strings.Split(name, "_")
	if len(s) >= 3 && s[2] == "table" {
		if s[0] == "create" {
			m.tp = createType
		} else if s[0] == "alter" {
			m.tp = alterType
		}
		m.tableName = s[1]
	}
	return nil
}

// Create create migration
func (m *Migrator) Create() error {
	if m.upName == "" || m.downName == "" {
		return errors.New("migration name required")
	}
	if err := m.createFile("up"); err != nil {
		return err
	}
	fmt.Printf("created migration %s\n", m.upName)
	if err := m.createFile("down"); err != nil {
		return err
	}
	fmt.Printf("created migration %s\n", m.downName)
	return nil
}

// Up up migration
func (m *Migrator) Up() error {
	if m.upName == "" || m.downName == "" {
		return errors.New("migration name required")
	}
	//get all migration files
	ms, err := m.getMigrations("up")
	if err != nil {
		return err
	}
	//get migration records
	rs, err := m.getMigrationRecords()
	if err != nil {
		return err
	}
	//except migrations
	for k, v := range ms {
		tmp := filepath.Base(v)
		for _, i := range rs {
			if tmp == i.migration+".up.sql" {
				ms = append(ms[:k], ms[k+1:]...)
			}
		}
	}
	if len(ms) == 0 {
		fmt.Println("everyting is up to date")
		return nil
	}
	//get sql，run
	var newBatch uint
	if len(rs) != 0 {
		newBatch = rs[0].batch + 1
	} else {
		newBatch = 1
	}
	for _, s := range ms {
		sql, err := ioutil.ReadFile(s)
		if err != nil {
			return err
		}
		tx, _ := m.db.Begin()
		_, err = tx.Exec(string(sql))
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		s = filepath.Base(s)
		s = strings.TrimRight(s, ".up.sql")
		_, err = tx.Exec("INSERT INTO migrations VALUES (null, ?, ?)", s, newBatch)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			return err
		}
		fmt.Printf("migrated %s success\n", s)
	}
	return nil
}

// Down down migration
func (m *Migrator) Down() error {
	if m.upName == "" || m.downName == "" {
		return errors.New("migration name required")
	}
	rs, err := m.getLastBatchRecords()
	if err != nil {
		return err
	}
	for _, v := range rs {
		fmt.Printf("%v\n", v)
		sql, err := ioutil.ReadFile(path.Join(m.path, v.migration) + ".down.sql")
		if err != nil {
			return err
		}
		tx, _ := m.db.Begin()
		// run file sql
		if _, err = tx.Exec(string(sql)); err != nil {
			_ = tx.Rollback()
			return err
		}
		// delete record
		_, err = tx.Exec("delete from migrations where id = ?", v.id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			return err
		}
		fmt.Printf("rollback %s success\n", v.migration)
	}
	return nil
}

func (m *Migrator) createFile(t string) error {
	name := m.upName
	if t == "down" {
		name = m.downName
	}
	f, err := os.Create(path.Join(m.path, name))
	if err != nil {
		return err
	}
	defer f.Close()

	sql := m.getSQL(t)
	if sql == "" {
		return nil
	}
	if _, err := f.WriteString(sql); err != nil {
		return err
	}
	return nil
}

func (m *Migrator) getMigrations(t string) ([]string, error) {
	matches, err := filepath.Glob(fmt.Sprintf("%s/*.%s.sql", m.path, t))
	if err != nil {
		return nil, err
	}
	return matches, nil
}

type record struct {
	id        uint
	migration string
	batch     uint
}

func (m *Migrator) getMigrationRecords() ([]record, error) {
	sql := "select id, migration, batch from migrations order by id desc"
	rows, err := m.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []record
	for rows.Next() {
		var r record
		if err := rows.Scan(&r.id, &r.migration, &r.batch); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func (m *Migrator) getLastBatchRecords() ([]record, error) {
	rs, err := m.getMigrationRecords()
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 {
		return nil, nil
	}
	var res []record
	lastBatch := rs[0].batch
	for _, v := range rs {
		if v.batch == lastBatch {
			res = append(res, v)
		}
	}
	return res, nil
}

func (m *Migrator) createMigrationTable() error {
	sqlContent := `CREATE TABLE IF NOT EXISTS migrations (
   id INT UNSIGNED AUTO_INCREMENT,
   migration VARCHAR(255),
   batch INT UNSIGNED,
   PRIMARY KEY (id)
)
`
	if _, err := m.db.Exec(sqlContent); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) getSQL(t string) string {
	sql := ""
	if m.tp == createType {
		if t == "up" {
			sql = m.createSQL()
		} else if t == "down" {
			sql = m.dropSQL()
		}
	}
	if m.tp == alterType {
		sql = m.alterSQL()
	}
	if m.tp != 0 {
		sql = strings.Replace(sql, "tb_name", m.tableName, 1)
	}
	return sql
}

func (m *Migrator) createSQL() string {
	return `CREATE TABLE IF NOT EXISTS tb_name (
   id INT UNSIGNED AUTO_INCREMENT,
   PRIMARY KEY (id)
)
`
}

func (m *Migrator) alterSQL() string {
	return "ALTER TABLE tb_name"
}

func (m *Migrator) dropSQL() string {
	return "DROP TABLE IF EXISTS tb_name"
}
