package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = stage(waitDone(out, done))
	}

	return out
}

func waitDone(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for v := range in {
			select {
			case <-done:
				continue
			default:
			}

			select {
			case out <- v:
			case <-done:
			}
		}
	}()

	return out
}
