// Package dbstorage реализует механизи взаимодействия с БД Postgresql
package dbstorage

import "database/sql"

const (
	restoreDBscript string = `CREATE TABLE IF NOT EXISTS metrics ( 
            name text, 
            type varchar(10), 
            value double precision,
            delta bigint,
 UNIQUE (name,type)
        );

CREATE OR REPLACE FUNCTION get() 
    RETURNS TABLE ( name text, 
            type varchar(10), 
            value double precision,
            delta bigint) AS $$
    SELECT name,type,value,delta FROM metrics
$$ LANGUAGE SQL STABLE;`
	SetMetricQuery string = `INSERT INTO metrics (name,type,value,delta)
        VALUES (@name,@type,@value,@delta)`
	GetMetricsQuery string = "SELECT * FROM get()"
)

type dbMetric struct {
	Name  sql.NullString
	MType sql.NullString
	Delta sql.NullInt64
	Value sql.NullFloat64
}
