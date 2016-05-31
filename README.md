# go2

A lightweight DC/OS service discovery provider. Looks up the IP and port given a certain distributed PID (Marathon app ID):

    $DPID -> http://$IP:$PORT

For example:

    /abc/def -> http://10.0.1.162:8652

To launch `go2` manually do the following:

    $ ./go2 &>/dev/null &

Once `go2` is running you can use it as so (assuming there's a Marathon app with the ID `/test` running):

    $ curl $(curl -s http://localhost:6969/?dpid=test)

In above command, the inner `curl` does the actual service discovery, returning something like `http://10.0.2.161:7192` which is then fed into the second `curl` command.

In production, you should launch via DC/OS Marathon:

    $ dcos marathon app add go2.json
