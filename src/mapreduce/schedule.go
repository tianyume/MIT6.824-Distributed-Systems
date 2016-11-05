package mapreduce

import "fmt"
import "sync"

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	var wg sync.WaitGroup
	for i := 0; i < ntasks; i++ {
		wg.Add(1)
		fmt.Printf("i: %d\n", i)
		fmt.Printf("len: %d\n", len(mr.files))
		go func(mr *Master, phase jobPhase, i int, nios int) {
			defer wg.Done()
			worker := <-mr.registerChannel
			doTaskArgs := DoTaskArgs{mr.jobName, mr.files[i], phase, i, nios}
			call(worker, "Worker.DoTask", doTaskArgs, new(struct{}))
			// Must use go func, or will be deadlocked
			go func(mr *Master, worker string) {
				mr.registerChannel <- worker
			}(mr, worker)
		}(mr, phase, i, nios)
	}
	wg.Wait()

	fmt.Printf("Schedule: %v phase done\n", phase)
}
