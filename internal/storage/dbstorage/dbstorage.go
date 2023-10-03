// Package dbstorage реализует механизи взаимодействия с БД Postgresql.
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

CREATE OR REPLACE FUNCTION save(_name TEXT,_type VARCHAR(10),_value DOUBLE PRECISION,_delta bigint)
 RETURNS void AS '
BEGIN
 INSERT INTO metrics (name,type,value,delta)
 VALUES (_name,_type,_value,_delta)
 ON CONFLICT(name,type) DO UPDATE
 SET value = EXCLUDED.value, delta = EXCLUDED.delta;
END;
' LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get() 
    RETURNS TABLE ( name text, 
            type varchar(10), 
            value double precision,
            delta bigint) AS $$
    SELECT name,type,value,delta FROM metrics
$$ LANGUAGE SQL STABLE;`
	SetMetricQuery  string = "SELECT save(@name,@type,@value,@delta)"
	GetMetricsQuery string = "SELECT * FROM get()"
)

type dbMetric struct {
	Name  sql.NullString
	MType sql.NullString
	Delta sql.NullInt64
	Value sql.NullFloat64
}
