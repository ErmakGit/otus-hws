package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out
	for k, stage := range stages {
		if k == 0 {
			out = stage(executeStage(in, done))
		} else {
			out = stage(executeStage(out, done))
		}
	}

	return out
}

func executeStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()

	return out
}
