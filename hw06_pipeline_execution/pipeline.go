package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// входной и выходной каналы для итерации по текущему стейджу
	var (
		stgIn  In
		stgOut Out
	)
	stgIn = in
	stgOut = stgIn

	if in == nil {
		return stgOut
	}

	for _, f := range stages {
		// для возможности прерывания всех каналов и ввиду того, что сигнатуру стейджей менять нельзя
		// в отдельных рутинах перепаковываем входной канал для каждого стейджа
		// в локальный входной канал. с возможностью терминации по done каналу
		lclIn := make(Bi)
		go func(stgIn In) {
			defer close(lclIn)
			for {
				select {
				case data, ok := <-stgIn:
					if !ok {
						return
					}
					lclIn <- data
				case <-done:
					return
				}
			}
		}(stgIn)

		// запускаем стейдж с локальным входным каналом
		stgOut = f(lclIn)
		// выходной канал текущего стейджа - это входной канал для следующего стейджа
		stgIn = stgOut
	}

	// выходной канал крайнего стейджа - целевой выходной канал обработчика
	return stgOut
}
