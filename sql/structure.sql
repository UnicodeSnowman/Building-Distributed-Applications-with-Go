CREATE DATABASE pluralsight_distributed_go;

CREATE TABLE sensor_readings (
  id integer NOT NULL,
  value double precision NOT NULL,
  sensor_id character varying NOT NULL,
  taken_on timestamp without time zone NOT NULL
);

CREATE SEQUENCE sensor_readings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE sensor_readings_id_seq OWNED BY sensor_readings.id;

ALTER TABLE ONLY sensor_readings ALTER COLUMN id SET DEFAULT nextval('sensor_readings_id_seq'::regclass);

CREATE TABLE sensors (
  id integer NOT NULL,
  name character varying NOT NULL
);

CREATE SEQUENCE sensors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE sensors_id_seq OWNED BY sensors.id;

ALTER TABLE ONLY sensors ALTER COLUMN id SET DEFAULT nextval('sensors_id_seq'::regclass);
