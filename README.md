# go2

go2 (pronounced: goto) is a lightweight DC/OS service discovery provider written in Go. Given a distributed PID (Marathon app ID) such as `/abc/def` it looks up an IP and port of a Mesos task running an instance of the Marathon app.

    $DPID -> http://$IP:$PORT

For example:

    /abc/def -> http://10.0.1.162:8652

## Development and testing

To launch `go2` manually do the following:

    $ ./go2 &>/dev/null &

Once `go2` is running you can use it as so (assuming there's a Marathon app with the ID `/test` running):

    $ curl $(curl -s http://localhost:6969/?dpid=test)

In above command, the inner `curl` does the actual service discovery, returning something like `http://10.0.2.161:7192` which is then fed into the second `curl` command. Note that you can also use escaped DPIDs such as `?dpid=%2Ftest%2Ft0`, especially useful when calling it from within a web browser.

## Production

In production, you should launch via DC/OS Marathon:

    $ dcos marathon app add go2.json

Then, `go2` is available on port `6969` of the public agent.

