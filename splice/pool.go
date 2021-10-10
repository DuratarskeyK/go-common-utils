package splice

import "sync"

var splicerPool = &sync.Pool{}

func GetCopier() *Copier {
	splicer := splicerPool.Get()
	if splicer != nil {
		return splicer.(*Copier)
	}

	return &Copier{}
}

func PutCopier(splicer *Copier) {
	splicer.doneInternal, splicer.Done = nil, nil
	splicer.once = nil
	splicer.returnError = nil

	splicerPool.Put(splicer)
}
