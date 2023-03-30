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

CREATE OR REPLACE FUNCTION save(TEXT, VARCHAR(10),DOUBLE PRECISION, bigint)
 RETURNS void AS '
BEGIN
 INSERT INTO metrics (name,type,value,delta)
 VALUES ($1,$2,$3,$4)
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
	SetMetricQuery  string = "SELECT * FROM save(@name,@type,@value,@delta)"
	GetMetricsQuery string = "SELECT * FROM get()"
)

type dbMetric struct {
	Name  sql.NullString
	MType sql.NullString
	Delta sql.NullInt64
	Value sql.NullFloat64
}
