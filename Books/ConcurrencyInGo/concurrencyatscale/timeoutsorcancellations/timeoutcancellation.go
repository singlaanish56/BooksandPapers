package timeoutsorcancellations


//now in this case the genration and the resultstream are perfectly handled by the done channels
//so they cancelled properly
//but the issue lies with the func call reallylongCalculation
//it could take a long time to cancel
func cancelTheLongCancellation(){
	var value interface{}

	done := make(chan interface{})
	valuestream := make(chan interface{})
	resultStream := make(chan interface{})

	select{
	case <-done:
		return
	case value = <-valuestream:
	}

	result := reallyLongCalculation(done, value)

	select{
	case <-done:
		return
	case resultStream<-result:
	}
}

//the best way to do this is to break the code into really small atomic functions
//: define the period within which our concurrent process is preemptable, and
//ensure that any functionality that takes more time than this period is itself preemptable.

// the best way to do this is to introduce the done channel in the function calls as well

func reallyLongCalculation(done <-chan interface{}, value interface{}) interface{} {
	intermediateResult := longCalculation(value)
	select{
	case <-done:
		return nil
	default:
	}

	return longCalculation(intermediateResult)
}

func longCalculation(value interface{}) interface{}{
	return value
}
