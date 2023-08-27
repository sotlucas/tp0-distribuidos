# Checks whether the server is running. If it is, it will print "OK".
docker run -it --network=tp0_testing_net --rm alpine /bin/sh -c "echo OK | nc server 12345"