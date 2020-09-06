package migrate

import (
	"database/sql"
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
func New(db *sql.DB, path, name string) *Migrator {
	migrator := &Migrator{
		db:        db,
		path:      path,
		tableName: name,
		ts:        time.Now().Format("20060102150405"),
	}
	migrator.upName = fmt.Sprintf("%s_%s.up.sql", migrator.ts, name)
	migrator.downName = fmt.Sprintf("%s_%s.down.sql", migrator.ts, name)
	s := strings.Split(name, "_")
	if len(s) == 3 && s[2] == "table" {
		if s[0] == "create" {
			migrator.tp = createType
		} else if s[0] == "alter" {
			migrator.tp = alterType
		}
		migrator.tableName = s[1]
	}
	if err := migrator.createMigrationTable(); err != nil {
		panic(err)
	}
	return migrator
}

// Create create migration
func (m *Migrator) Create() error {
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

// Up up migration
func (m *Migrator) Up() error {
	//获取migrations文件
	ms, err := m.getMigrations("up")
	if err != nil {
		return err
	}
	//排除已执行的migrations记录
	rs, err := m.getMigrationRecords()
	if err != nil {
		return err
	}
	//排除已执行的migrations
	for k, v := range ms {
		tmp := filepath.Base(v)
		for _, i := range rs {
			if tmp == i+".up.sql" {
				ms = append(ms[:k], ms[k+1:]...)
			}
		}
	}
	if len(ms) == 0 {
		fmt.Println("everyting is up to date")
		return nil
	}
	//获取sql，执行
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
		_, err = tx.Exec("INSERT INTO migrations VALUES (null, ?)", s)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		fmt.Printf("migrated %s success\n", s)
	}
	return nil
}

// Down down migration
func (m *Migrator) Down() error { return nil }

func (m *Migrator) getMigrations(t string) ([]string, error) {
	matches, err := filepath.Glob(fmt.Sprintf("%s/*.%s.sql", m.path, t))
	if err != nil {
		return nil, err
	}
	return matches, nil
}
func (m *Migrator) getMigrationRecords() ([]string, error) {
	sql := "select migration from migrations order by id desc"
	rows, err := m.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var i string
		if err := rows.Scan(&i); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func (m *Migrator) createMigrationTable() error {
	sqlContent := `CREATE TABLE IF NOT EXISTS migrations (
   id INT UNSIGNED AUTO_INCREMENT,
   migration VARCHAR(255),
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
