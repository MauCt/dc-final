Architecture Document
=====================
					                                            ______> Worker
	    	      MSG		            Jobs                |
Client --------> API ------> Controller -----------> Scheduler |-----> Worker
							                                     |_____> Worker
							                                         RPC