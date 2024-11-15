--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: dblink; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS dblink WITH SCHEMA public;


--
-- Name: EXTENSION dblink; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION dblink IS 'connect to other PostgreSQL databases from within a database';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: trend; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.trend (
    temperature text not null,
    humidity text not null,
    pressure text not null
);


--
-- Name: weather_trend(interval); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.weather_trend(since interval) RETURNS SETOF public.trend
    LANGUAGE plpgsql
    AS $$
begin
  return query (
    select 
      case when temperature > 0 then 'increasing' when temperature < 0 then 'decreasing' else 'stable' end temperature,
      case when humidity > 0 then 'increasing' when humidity < 0 then 'decreasing' else 'stable' end humidity,
      case when pressure > 0 then 'increasing' when pressure < 0 then 'decreasing' else 'stable' end pressure
    from (
      select
        regr_slope(temperature, extract(epoch from recorded_at)) as temperature,
        regr_slope(humidity, extract(epoch from recorded_at)) as humidity,
        regr_slope(pressure, extract(epoch from recorded_at)) as pressure
      from public.measurements m 
      where recorded_at > now() - since
    )
  );
end; $$;


--
-- Name: measurements; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.measurements (
    id bigint NOT NULL,
    recorded_at timestamp without time zone NOT NULL,
    temperature double precision NOT NULL,
    humidity double precision NOT NULL,
    pressure double precision NOT NULL,
    location character varying
);


--
-- Name: measurements_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.measurements_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: measurements_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.measurements_id_seq OWNED BY public.measurements.id;


--
-- Name: measurements id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.measurements ALTER COLUMN id SET DEFAULT nextval('public.measurements_id_seq'::regclass);


--
-- Name: measurements measurements_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.measurements
    ADD CONSTRAINT measurements_pk PRIMARY KEY (id);


--
-- Name: measurements_recorded_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX measurements_recorded_at_idx ON public.measurements USING btree (recorded_at DESC);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: -
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

