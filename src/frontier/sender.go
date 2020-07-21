package frontier

import (
	"sync"
)

var Sender = &sender{}

type job struct {
	conn    Conn
	message []byte
}

type sender struct {
	init sync.Once
	jobs chan *job
	pool sync.Pool
}

func (s *sender) Push(conn Conn, message []byte) {
	s.init.Do(func() {
		s.jobs = make(chan *job, 2)
		s.pool.New = func() interface{} {
			return new(job)
		}
		go func() {
			for job := range s.jobs {
				m := job.message
				job.conn.Writer(m)
				s.pool.Put(job)
			}
		}()
	})

	j := s.pool.Get().(*job)
	j.conn, j.message = conn, message
	s.jobs <- j
}
