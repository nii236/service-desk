package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bgentry/que-go"
	"github.com/jackc/pgx"
)

type printNameArgs struct {
	Name string
}

func QueueTest() {
	printName := func(j *que.Job) error {
		var args printNameArgs
		if err := json.Unmarshal(j.Args, &args); err != nil {
			return err
		}
		fmt.Printf("Hello %s!\n", args.Name)
		return nil
	}

	sendEmail := func(j *que.Job) error { return nil }
	sendComment := func(j *que.Job) error { return nil }

	pgxcfg, err := pgx.ParseURI("postgres://dev:dev@localhost:5432/dev")
	if err != nil {
		log.Fatal(err)
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgxcfg,
		AfterConnect: func(conn *pgx.Conn) error {
			_, err = conn.Exec(`
			CREATE TABLE IF NOT EXISTS que_jobs
			(
			  priority    smallint    NOT NULL DEFAULT 100,
			  run_at      timestamptz NOT NULL DEFAULT now(),
			  job_id      bigserial   NOT NULL,
			  job_class   text        NOT NULL,
			  args        json        NOT NULL DEFAULT '[]'::json,
			  error_count integer     NOT NULL DEFAULT 0,
			  last_error  text,
			  queue       text        NOT NULL DEFAULT '',
			
			  CONSTRAINT que_jobs_pkey PRIMARY KEY (queue, priority, run_at, job_id)
			);
			
			COMMENT ON TABLE que_jobs IS '3';
			`)
			err = que.PrepareStatements(conn)
			return err
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	defer pgxpool.Close()

	if err != nil {
		log.Fatalln(err)
	}
	qc := que.NewClient(pgxpool)
	wm := que.WorkMap{
		"PrintName":   printName,
		"SendEmail":   sendEmail,
		"SendComment": sendComment,
	}
	workers := que.NewWorkerPool(qc, wm, 2) // create a pool w/ 2 workers
	go workers.Start()                      // work jobs in another goroutine

	args, err := json.Marshal(printNameArgs{Name: "bgentry"})
	if err != nil {
		log.Fatal(err)
	}

	j1 := &que.Job{
		Type: "PrintName",
		Args: args,
	}
	if err := qc.Enqueue(j1); err != nil {
		log.Fatal(err)
	}

	j2 := &que.Job{
		Type:  "PrintName",
		RunAt: time.Now().UTC().Add(5 * time.Second), // delay 30 seconds
		Args:  args,
	}
	if err := qc.Enqueue(j2); err != nil {
		log.Fatal(err)
	}

	time.Sleep(35 * time.Second) // wait for while

	workers.Shutdown()
}
