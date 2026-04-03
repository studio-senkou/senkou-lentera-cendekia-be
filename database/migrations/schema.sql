\restrict VrNAwHwvLAm7hfW7oJgiAzhJpNcGdMRt9KuYYggFdbuwO4dpoGx7IAam1KAsf4Q

-- Dumped from database version 17.7
-- Dumped by pg_dump version 17.9 (Ubuntu 17.9-0ubuntu0.25.10.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: btree_gin; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS btree_gin WITH SCHEMA public;


--
-- Name: EXTENSION btree_gin; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION btree_gin IS 'support for indexing common datatypes in GIN';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: blogs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blogs (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    content text NOT NULL,
    author_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: blogs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blogs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blogs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blogs_id_seq OWNED BY public.blogs.id;


--
-- Name: classes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.classes (
    id uuid NOT NULL,
    classname character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: meeting_session_proofs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.meeting_session_proofs (
    id integer NOT NULL,
    meeting_id integer NOT NULL,
    student_proof character varying(255),
    student_signature character varying(255),
    mentor_proof character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: meeting_session_proofs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.meeting_session_proofs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: meeting_session_proofs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.meeting_session_proofs_id_seq OWNED BY public.meeting_session_proofs.id;


--
-- Name: meeting_sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.meeting_sessions (
    id integer NOT NULL,
    student_id integer NOT NULL,
    session_date date NOT NULL,
    session_time time without time zone NOT NULL,
    duration_minutes smallint NOT NULL,
    status character varying(20) NOT NULL,
    note text,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    mentor_id integer NOT NULL
);


--
-- Name: meeting_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.meeting_sessions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: meeting_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.meeting_sessions_id_seq OWNED BY public.meeting_sessions.id;


--
-- Name: mentors; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.mentors (
    id integer NOT NULL,
    user_id integer NOT NULL,
    class_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: mentors_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.mentors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: mentors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.mentors_id_seq OWNED BY public.mentors.id;


--
-- Name: quiz_answers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.quiz_answers (
    id integer NOT NULL,
    attempt_id integer NOT NULL,
    question_id integer NOT NULL,
    option_id integer NOT NULL,
    is_correct boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: quiz_answers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.quiz_answers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: quiz_answers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.quiz_answers_id_seq OWNED BY public.quiz_answers.id;


--
-- Name: quiz_attempts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.quiz_attempts (
    id integer NOT NULL,
    quiz_id integer NOT NULL,
    user_id integer NOT NULL,
    status character varying(20) DEFAULT 'in_progress'::character varying NOT NULL,
    score numeric(5,2),
    started_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    submitted_at timestamp without time zone,
    reset_at timestamp without time zone,
    reset_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    question_ids integer[],
    option_order jsonb,
    current_question_index integer DEFAULT 0
);


--
-- Name: quiz_attempts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.quiz_attempts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: quiz_attempts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.quiz_attempts_id_seq OWNED BY public.quiz_attempts.id;


--
-- Name: quiz_options; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.quiz_options (
    id integer NOT NULL,
    question_id integer NOT NULL,
    option_text text NOT NULL,
    is_correct boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: quiz_options_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.quiz_options_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: quiz_options_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.quiz_options_id_seq OWNED BY public.quiz_options.id;


--
-- Name: quiz_questions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.quiz_questions (
    id integer NOT NULL,
    quiz_id integer NOT NULL,
    question_text text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: quiz_questions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.quiz_questions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: quiz_questions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.quiz_questions_id_seq OWNED BY public.quiz_questions.id;


--
-- Name: quiz_quizzes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.quiz_quizzes (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    passing_score smallint DEFAULT 70 NOT NULL,
    time_limit_minutes smallint,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: quiz_quizzes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.quiz_quizzes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: quiz_quizzes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.quiz_quizzes_id_seq OWNED BY public.quiz_quizzes.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: static_assets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.static_assets (
    id integer NOT NULL,
    asset_name character varying(255) NOT NULL,
    asset_type character varying(50) NOT NULL,
    asset_url text NOT NULL,
    asset_description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: static_assets_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.static_assets_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: static_assets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.static_assets_id_seq OWNED BY public.static_assets.id;


--
-- Name: student_plans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.student_plans (
    id integer NOT NULL,
    student_id integer NOT NULL,
    total_sessions integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: student_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.student_plans_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: student_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.student_plans_id_seq OWNED BY public.student_plans.id;


--
-- Name: students; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.students (
    id integer NOT NULL,
    user_id integer NOT NULL,
    class_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: students_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.students_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: students_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.students_id_seq OWNED BY public.students.id;


--
-- Name: testimonials; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.testimonials (
    id integer NOT NULL,
    testimoner_name character varying(255) NOT NULL,
    testimoner_current_position character varying(255),
    testimoner_previous_position character varying(255),
    testimoner_photo text,
    testimony_text text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: testimonials_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.testimonials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: testimonials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.testimonials_id_seq OWNED BY public.testimonials.id;


--
-- Name: user_has_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_has_tokens (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: user_has_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_has_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_has_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_has_tokens_id_seq OWNED BY public.user_has_tokens.id;


--
-- Name: user_has_tokens_user_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_has_tokens_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_has_tokens_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_has_tokens_user_id_seq OWNED BY public.user_has_tokens.user_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    password text NOT NULL,
    role character varying(20) DEFAULT 'user'::character varying NOT NULL,
    email_verified_at timestamp without time zone,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: blogs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blogs ALTER COLUMN id SET DEFAULT nextval('public.blogs_id_seq'::regclass);


--
-- Name: meeting_session_proofs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meeting_session_proofs ALTER COLUMN id SET DEFAULT nextval('public.meeting_session_proofs_id_seq'::regclass);


--
-- Name: meeting_sessions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meeting_sessions ALTER COLUMN id SET DEFAULT nextval('public.meeting_sessions_id_seq'::regclass);


--
-- Name: mentors id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mentors ALTER COLUMN id SET DEFAULT nextval('public.mentors_id_seq'::regclass);


--
-- Name: quiz_answers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_answers ALTER COLUMN id SET DEFAULT nextval('public.quiz_answers_id_seq'::regclass);


--
-- Name: quiz_attempts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_attempts ALTER COLUMN id SET DEFAULT nextval('public.quiz_attempts_id_seq'::regclass);


--
-- Name: quiz_options id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_options ALTER COLUMN id SET DEFAULT nextval('public.quiz_options_id_seq'::regclass);


--
-- Name: quiz_questions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_questions ALTER COLUMN id SET DEFAULT nextval('public.quiz_questions_id_seq'::regclass);


--
-- Name: quiz_quizzes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_quizzes ALTER COLUMN id SET DEFAULT nextval('public.quiz_quizzes_id_seq'::regclass);


--
-- Name: static_assets id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.static_assets ALTER COLUMN id SET DEFAULT nextval('public.static_assets_id_seq'::regclass);


--
-- Name: student_plans id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.student_plans ALTER COLUMN id SET DEFAULT nextval('public.student_plans_id_seq'::regclass);


--
-- Name: students id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.students ALTER COLUMN id SET DEFAULT nextval('public.students_id_seq'::regclass);


--
-- Name: testimonials id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.testimonials ALTER COLUMN id SET DEFAULT nextval('public.testimonials_id_seq'::regclass);


--
-- Name: user_has_tokens id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_has_tokens ALTER COLUMN id SET DEFAULT nextval('public.user_has_tokens_id_seq'::regclass);


--
-- Name: user_has_tokens user_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_has_tokens ALTER COLUMN user_id SET DEFAULT nextval('public.user_has_tokens_user_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: blogs blogs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blogs
    ADD CONSTRAINT blogs_pkey PRIMARY KEY (id);


--
-- Name: classes classes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.classes
    ADD CONSTRAINT classes_pkey PRIMARY KEY (id);


--
-- Name: meeting_session_proofs meeting_session_proofs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meeting_session_proofs
    ADD CONSTRAINT meeting_session_proofs_pkey PRIMARY KEY (id);


--
-- Name: meeting_sessions meeting_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meeting_sessions
    ADD CONSTRAINT meeting_sessions_pkey PRIMARY KEY (id);


--
-- Name: mentors mentors_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mentors
    ADD CONSTRAINT mentors_pkey PRIMARY KEY (id);


--
-- Name: quiz_answers quiz_answers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_answers
    ADD CONSTRAINT quiz_answers_pkey PRIMARY KEY (id);


--
-- Name: quiz_attempts quiz_attempts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_attempts
    ADD CONSTRAINT quiz_attempts_pkey PRIMARY KEY (id);


--
-- Name: quiz_options quiz_options_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_options
    ADD CONSTRAINT quiz_options_pkey PRIMARY KEY (id);


--
-- Name: quiz_questions quiz_questions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_questions
    ADD CONSTRAINT quiz_questions_pkey PRIMARY KEY (id);


--
-- Name: quiz_quizzes quiz_quizzes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_quizzes
    ADD CONSTRAINT quiz_quizzes_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: static_assets static_assets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.static_assets
    ADD CONSTRAINT static_assets_pkey PRIMARY KEY (id);


--
-- Name: student_plans student_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.student_plans
    ADD CONSTRAINT student_plans_pkey PRIMARY KEY (id);


--
-- Name: students students_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_pkey PRIMARY KEY (id);


--
-- Name: testimonials testimonials_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.testimonials
    ADD CONSTRAINT testimonials_pkey PRIMARY KEY (id);


--
-- Name: quiz_answers uq_quiz_answers_attempt_question; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_answers
    ADD CONSTRAINT uq_quiz_answers_attempt_question UNIQUE (attempt_id, question_id);


--
-- Name: user_has_tokens user_has_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_has_tokens
    ADD CONSTRAINT user_has_tokens_pkey PRIMARY KEY (id);


--
-- Name: user_has_tokens user_has_tokens_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_has_tokens
    ADD CONSTRAINT user_has_tokens_user_id_key UNIQUE (user_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_blogs_author_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_blogs_author_id ON public.blogs USING btree (author_id);


--
-- Name: idx_meeting_session_proofs_meeting_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_meeting_session_proofs_meeting_id ON public.meeting_session_proofs USING btree (meeting_id);


--
-- Name: idx_meeting_sessions_mentor_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_meeting_sessions_mentor_id ON public.meeting_sessions USING btree (mentor_id);


--
-- Name: idx_meeting_sessions_student_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_meeting_sessions_student_id ON public.meeting_sessions USING btree (student_id);


--
-- Name: idx_mentors_class_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_mentors_class_id ON public.mentors USING btree (class_id);


--
-- Name: idx_mentors_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_mentors_user_id ON public.mentors USING btree (user_id);


--
-- Name: idx_quiz_answers_attempt_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_quiz_answers_attempt_id ON public.quiz_answers USING btree (attempt_id);


--
-- Name: idx_quiz_attempts_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_quiz_attempts_status ON public.quiz_attempts USING btree (status);


--
-- Name: idx_quiz_attempts_user_quiz; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_quiz_attempts_user_quiz ON public.quiz_attempts USING btree (user_id, quiz_id);


--
-- Name: idx_quiz_options_question_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_quiz_options_question_id ON public.quiz_options USING btree (question_id);


--
-- Name: idx_quiz_questions_quiz_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_quiz_questions_quiz_id ON public.quiz_questions USING btree (quiz_id);


--
-- Name: idx_quiz_quizzes_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_quiz_quizzes_is_active ON public.quiz_quizzes USING btree (is_active);


--
-- Name: idx_static_assets_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_static_assets_url ON public.static_assets USING btree (asset_url);


--
-- Name: idx_student_plans_student_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_student_plans_student_id ON public.student_plans USING btree (student_id);


--
-- Name: idx_students_class_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_students_class_id ON public.students USING btree (class_id);


--
-- Name: idx_students_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_students_user_id ON public.students USING btree (user_id);


--
-- Name: idx_user_tokens; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_tokens ON public.user_has_tokens USING btree (user_id);


--
-- Name: blogs fk_blogs_author; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blogs
    ADD CONSTRAINT fk_blogs_author FOREIGN KEY (author_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: meeting_sessions fk_mentor_user_mt_sessions; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meeting_sessions
    ADD CONSTRAINT fk_mentor_user_mt_sessions FOREIGN KEY (mentor_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: mentors fk_mentors_class; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mentors
    ADD CONSTRAINT fk_mentors_class FOREIGN KEY (class_id) REFERENCES public.classes(id) ON DELETE CASCADE;


--
-- Name: mentors fk_mentors_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mentors
    ADD CONSTRAINT fk_mentors_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: meeting_session_proofs fk_mt_session_proof; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meeting_session_proofs
    ADD CONSTRAINT fk_mt_session_proof FOREIGN KEY (meeting_id) REFERENCES public.meeting_sessions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: quiz_answers fk_quiz_answers_attempt_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_answers
    ADD CONSTRAINT fk_quiz_answers_attempt_id FOREIGN KEY (attempt_id) REFERENCES public.quiz_attempts(id) ON DELETE CASCADE;


--
-- Name: quiz_answers fk_quiz_answers_option_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_answers
    ADD CONSTRAINT fk_quiz_answers_option_id FOREIGN KEY (option_id) REFERENCES public.quiz_options(id) ON DELETE CASCADE;


--
-- Name: quiz_answers fk_quiz_answers_question_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_answers
    ADD CONSTRAINT fk_quiz_answers_question_id FOREIGN KEY (question_id) REFERENCES public.quiz_questions(id) ON DELETE CASCADE;


--
-- Name: quiz_attempts fk_quiz_attempts_quiz_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_attempts
    ADD CONSTRAINT fk_quiz_attempts_quiz_id FOREIGN KEY (quiz_id) REFERENCES public.quiz_quizzes(id) ON DELETE CASCADE;


--
-- Name: quiz_attempts fk_quiz_attempts_reset_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_attempts
    ADD CONSTRAINT fk_quiz_attempts_reset_by FOREIGN KEY (reset_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: quiz_attempts fk_quiz_attempts_user_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_attempts
    ADD CONSTRAINT fk_quiz_attempts_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: quiz_options fk_quiz_options_question_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_options
    ADD CONSTRAINT fk_quiz_options_question_id FOREIGN KEY (question_id) REFERENCES public.quiz_questions(id) ON DELETE CASCADE;


--
-- Name: quiz_questions fk_quiz_questions_quiz_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quiz_questions
    ADD CONSTRAINT fk_quiz_questions_quiz_id FOREIGN KEY (quiz_id) REFERENCES public.quiz_quizzes(id) ON DELETE CASCADE;


--
-- Name: meeting_sessions fk_student_mt_sessions; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.meeting_sessions
    ADD CONSTRAINT fk_student_mt_sessions FOREIGN KEY (student_id) REFERENCES public.students(id) ON DELETE CASCADE;


--
-- Name: student_plans fk_student_plans; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.student_plans
    ADD CONSTRAINT fk_student_plans FOREIGN KEY (student_id) REFERENCES public.students(id) ON DELETE CASCADE;


--
-- Name: students fk_students_class; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT fk_students_class FOREIGN KEY (class_id) REFERENCES public.classes(id) ON DELETE CASCADE;


--
-- Name: students fk_students_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT fk_students_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: user_has_tokens fk_user_tokens; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_has_tokens
    ADD CONSTRAINT fk_user_tokens FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict VrNAwHwvLAm7hfW7oJgiAzhJpNcGdMRt9KuYYggFdbuwO4dpoGx7IAam1KAsf4Q


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250920160635'),
    ('20250920160830'),
    ('20250920161002'),
    ('20250920161324'),
    ('20250920161834'),
    ('20250920162001'),
    ('20250920162804'),
    ('20250920163459'),
    ('20250920233136'),
    ('20250920233216'),
    ('20250920235530'),
    ('20260103180100'),
    ('20260403000001'),
    ('20260403000002'),
    ('20260403000003'),
    ('20260403000004'),
    ('20260403000005'),
    ('20260403144400'),
    ('20260403145100'),
    ('20260403160200');
