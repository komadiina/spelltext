before stress testing, log into the application using credentials `{username=oggnjen, password=any}`, and then start the stress test:
```sh
$ py stresstest.py --help
usage: stresstest.py [-h] [-n N_CALLS] [-i INTERVAL] [-c CONNECT_TIMEOUT]

options:
  -h, --help            show this help message and exit
  -n, --n-calls N_CALLS
                        number of calls per method. default=50
  -i, --interval INTERVAL
                        interval between calls, in milliseconds. default=0.01
  -c, --connect-timeout CONNECT_TIMEOUT
                        maximum connect timeout/wait-per-response, in seconds. default=1

$ py stresstest.py -n 50 # invokes each grpc method 50 times
$ py stresstest.py -i 50 #
```