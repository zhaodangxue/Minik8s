package internal

import (
	"sync"
)

var jobs = &sync.Map{}

func init(){
}

func Jobs() *sync.Map{
	return jobs
}

func ReplaceJobs(newJobs *sync.Map){
	jobs = newJobs
}
