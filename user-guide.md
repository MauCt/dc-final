User Guide
==========
Welcome to this distributed system. To use it, you can follow the following steps

- First you have to switch to the 'dc-final' folder, then turn off the GOMODULE with this command:
	
	$ export GO111MODULE=off

  (be sure to type this command in every new terminal you create)

- Then run the program with the following command 
	
	$ go run main.go

- Now that the program is running, open a second bash terminal (the following instruction will be done
  in this terminal. 
- Now switch to the 'worker' folder and create the worker 
	
	$ go run main.go --controller tcp://localhost:40899 --worker-name worker1

- Open a third bash terminal (all of the following commands will be done in this terminal)
  and login to the system with the following command:

	$ curl -u username:password -X POST http://localhost:8080/login

  (this command returns a token. In all of the following commands, you must replace <ACCESS_TOKEN>
  with this token)

- Once you are logged in and the program is running, you can check the info of the server
  and the active workloads with this command

	$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status


- Create a workload

	$ curl -X POST -H "Authorization: Bearer <ACCESS_TOKEN>" --form 'filter="FILTER"' http://localhost:8080/workloads

  (be sure to replace FILTER with the filter you want to use)

- To see the workload information, type the following command and replace workload_id with the 
  workload ID provided when you created the workload

	$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/workloads/<workload_id>

- To upload an image to a workload use this command and replace image_name with the name of
  your image and worker_id with the id of the worker

	$python stress_test.py -action push -workload-id <workload_id> -token <token> -frames-path frames

- You can check the workload again to see the changes with the following command

	$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/workloads/<workload_id>

- Once everything is done, you can logout of the system

	$ curl -X DELETE -H "Authorization: <ACCESS_TOKEN>" http://localhost:8080/logout

	

	