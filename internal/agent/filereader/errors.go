package filereader

import "errors"

var ErrorNoFileName = errors.New("no file name provided")

var ErrorEmptyInputParams = errors.New("file and input arguments are both empty")
var ErrorWrongInputType = errors.New("binary input from argument is not supported")
