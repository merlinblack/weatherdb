--
-- PostgreSQL database dump
--

--
-- Name: trend; Type: TABLE; Schema: public
--

CREATE TABLE public.trend (
    temperature text NOT NULL,
    humidity text NOT NULL,
    pressure text NOT NULL
);


--
-- Name: weather_trend(interval); Type: FUNCTION; Schema: public;
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
-- Name: measurements; Type: TABLE; Schema: public; 
--

CREATE TABLE public.measurements (
    id bigint NOT NULL,
    recorded_at timestamp without time zone NOT NULL,
    temperature double precision NOT NULL,
    humidity double precision NOT NULL,
    pressure double precision NOT NULL,
    location character varying,
    hour bigint
);


--
-- Name: measurements_id_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE public.measurements_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: measurements_id_seq; Type: SEQUENCE OWNED BY; Schema: public; 
--

ALTER SEQUENCE public.measurements_id_seq OWNED BY public.measurements.id;


--
-- Name: measurements id; Type: DEFAULT; Schema: public; 
--

ALTER TABLE ONLY public.measurements ALTER COLUMN id SET DEFAULT nextval('public.measurements_id_seq'::regclass);


--
-- Name: measurements measurements_pk; Type: CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public.measurements
    ADD CONSTRAINT measurements_pk PRIMARY KEY (id);


--
-- Name: measurements_recorded_at_idx; Type: INDEX; Schema: public; 
--

CREATE INDEX measurements_recorded_at_idx ON public.measurements USING btree (recorded_at DESC);


