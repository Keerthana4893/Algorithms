package main

import "log"

func main() {
	/* In this main fucntion I have used DI for the over use of init func and instances which can be used anywhere in the dependecy tree. Google has recently launched its DI containers called wire. Albeit, I have used dig library DI in this.*/
	req_proc, prop_proc := dependencyInjections()
	log.Printf("Got %v and %v", req_proc, prop_proc)
}

func dependencyInjections() (Req_proc, Prop_proc) {
	db, err := NewDBClient()
	if err != nil {
		panic("Couldn't get connected to DB")
	}

	reAlgo := NewReq_m_algo()
	proAlgo := NewProp_m_algo()
	req_proc := NewReq_proc(db, reAlgo)
	prop_proc := NewProp_proc(db, proAlgo)
	return req_proc, prop_proc
} 
