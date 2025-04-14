package errorpropogation

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

//now error propogation is important in a concurrent system, because the error might make sense
// in a particulat context but it doesnt make sense in the parent context.
// for example a stack trace might sense at low level, but its not a proper error to display to a client
// such errors need to be wrapped, in order  to make sense and for it to be a well formed error

// aspects of a well formed error
type MyError struct{
	Inner error
	Message string
	StackTrace string
	Misc map[string]interface{}
}

func wrapError(err error, messagedef string, msgArgs ...interface{}) MyError{
	return MyError{
		Inner: err,
		Message: fmt.Sprintf(messagedef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc: make(map[string]interface{}),
	}
}


func (err MyError) Error() string{
	return err.Message
}


//error at a low level

type LowLevelErr struct{
	error
}

func isGloballyExec(path string) (bool, error){
	info, err := os.Stat(path)
	if err!=nil{
		return false, LowLevelErr{(wrapError(err, err.Error()))}
	}

	return info.Mode().Perm()&0100 == 0100, nil
}

//internediate

type IntermediateErr struct{
	error
}

func runJob(id string) error{
	const path = "/bad/job/path"
	isExec, err := isGloballyExec(path)
	if err != nil{
		return err // because of this in the path we still get a error path which are not yet readable for the humans
	}else if isExec==false{
		return wrapError(nil, "job binary not executable")
	}

	return exec.Command(path,"--id="+id).Run()
}

//top level or client level function

func handleError(key int, err error, message string){
	log.SetPrefix(fmt.Sprintf("[logid: %v]", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key , message)
}

func RunErrorProp(){
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.LUTC)

	err := runJob("1")
	if err !=nil{
		msg:="there was an unexpected issue; please report this as a bug"
		if _, ok := err.(IntermediateErr); ok{
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}


//more improvemnt on the run job is to  convert the error thrown by the low level to an intermediate error

func runJob2(id string) error{
	const path = "/bad/job/path"
	isExec, err := isGloballyExec(path)
	if err != nil{
		return IntermediateErr{ wrapError(
			err,
			"cannot run job %q: requisite binaries not available",
			id,
		)}
	}else if isExec==false{
		return wrapError(nil, "cannot run job %q: requisite binaries not available", id)
	}

	return exec.Command(path,"--id="+id).Run()
}

func RunErrorPropBetter(){
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.LUTC)

	err := runJob2("1")
	if err !=nil{
		msg:="there was an unexpected issue; please report this as a bug"
		if _, ok := err.(IntermediateErr); ok{
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
