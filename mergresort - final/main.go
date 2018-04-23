package main

import (
	// Import the entire framework for interracting with SDAccel from Go (including bundled verilog)
	_ "github.com/ReconfigureIO/sdaccel"

	// Use the new AXI protocol package for interracting with memory
	aximemory "github.com/ReconfigureIO/sdaccel/axi/memory"
	axiprotocol "github.com/ReconfigureIO/sdaccel/axi/protocol"
)



func replaceItem(m <-chan uint64,result chan<- uint64,item int,size int,input uint64)  {
	go func() {
		replace:=make(chan uint64)
		for i:=0;i<item;i++ {
			replace<-<-m
		}

		replace<-input
		<-m
		for i:=item;i<size-1;i++ {
			replace<-<-m
		}

		for i:=0;i<size ;i++  {
			result<-<-replace
		}

	}()
}
func getItem(input <-chan uint64,result chan<- uint64,size int,item int)(uint64) {
	swap:=make(chan uint64)
	go func() {

		replace:=make(chan uint64)

		for i:=0;i<item ;i++  {
			replace<-<-input
		}
		swap<-<-input
		replace<-<-swap

		for i:=item+1;i<size ;i++  {
			replace<-<-input
		}

		for i:=0;i<size ;i++  {
			result<-<-replace
		}
	}()

	return <-swap
}

func mergesort_iterative(input <-chan uint64,result chan<- uint64,size int,item int){
	go func() {
		temparr:= make(chan uint64)
		replacearr:= make(chan uint64)
		for i:=0;i<20;i++{
			replacearr<-<-input

		}
		right:=0
		rend:=0
		i:=0
		j:=0
		m:=0
		right++
		rend++
		i++
		j++
		m++
		for k:= 1; k < size; k *= 2 {
			//at each partition size, sort and merge
			for  left := 0; left + k < size; left += k*2 {
				//store the start of the right partition and its end
				right = left + k
				rend = right + k

				//if the partitions are uneven, readjust the end
				if rend > size{
					rend = size
				}
				m = left
				i = left
				j = right

				//merge
				for i < right && j < rend {


					if getItem(replacearr,replacearr,size,i) <=getItem(replacearr,replacearr,size,j) {
						replaceItem(temparr,temparr,m,size,getItem(replacearr,replacearr,size,i))
						//temparr[m] = arr[i]
						i++
					} else {
						replaceItem(temparr,temparr,m,size,getItem(replacearr,replacearr,size,j))
						//temparr[m] = arr[j]
						j++
					}
					m++
				}
				for i < right {
					//	temparr[m] = arr[i]
					replaceItem(temparr,temparr,m,size,getItem(replacearr,replacearr,size,i))
					i++
					m++
				}
				for j < rend {
					//temparr[m] = arr[j]
					replaceItem(temparr,temparr,m,size,getItem(replacearr,replacearr,size,j))
					j++
					m++
				}
				//copy from temp array into initial array
				for m = left; m < rend; m++ {
					replaceItem(replacearr,replacearr,m,size,getItem(temparr,temparr,size,m))
					//	arr[m] = temparr[m]
				}
			}
		}

	}()
}

func Top(
// Three operands from the host. Pointers to the input data and the space for the result in shared
// memory and the length of the input data so the FPGA knows what to expect.
	inputData uintptr,
	outputData uintptr,
	length uint32,

// Set up channels for interacting with the shared memory
	memReadAddr chan<- axiprotocol.Addr,
	memReadData <-chan axiprotocol.ReadData,

	memWriteAddr chan<- axiprotocol.Addr,
	memWriteData chan<- axiprotocol.WriteData,
	memWriteResp <-chan axiprotocol.WriteResp) {


	// Read all of the input data into a channel
	inputChan := make(chan uint64)
	//outputChan := make(chan uint64)
	go aximemory.ReadBurstUInt64(
		memReadAddr, memReadData, true, inputData, length, inputChan)


	go func() {
		//replaceItem(inputChan, inputChan, 0, 20, 1234)
		//mer(inputChan,inputChan,20)
		mergesort_iterative(inputChan,inputChan,20,0)
	}()


	// Write the results to shared memory
	aximemory.WriteBurstUInt64(
		memWriteAddr, memWriteData, memWriteResp, true, outputData, 20, inputChan)
}
