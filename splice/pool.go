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
	splicer.GlobalStop = nil
	splicer.Conn1ToConn2 = 0
	splicer.Conn2ToConn1 = 0

	splicer.stopChan = nil
	splicer.once = nil
	splicer.returnError = nil

	splicerPool.Put(splicer)
}
