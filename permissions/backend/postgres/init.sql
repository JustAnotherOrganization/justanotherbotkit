--
-- PostgreSQL database dump
--

-- Dumped from database version 10.1
-- Dumped by pg_dump version 10.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: permissions; Type: TABLE; Schema: public; Owner: ing
--

CREATE TABLE permissions (
    id character varying(255) NOT NULL,
    perms character varying(255)[] NOT NULL,
    groups integer[] NOT NULL
);


ALTER TABLE permissions OWNER TO ing;

--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: ing
--

ALTER TABLE ONLY permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

