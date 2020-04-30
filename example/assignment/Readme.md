# LIVE STREAMER

A simple client-server application to live stream the video captured via webcam/video file over the network.
We start a server echoing data on the first stream the client opens,then connect with a client. 
We then start capturing the video frames using GoCV module, serialize the image matrix and write the frame length and frame data to the stream in respective order.
Server first receives the length of the frame data in a fixed length container and then receives the complete frame.
It then coverts the received frame back to image mat and displays at the server side.

### Setup

1. Complete the installation of mpquic as mentioned int the project readme.

2. Install the GoCV package using the steps here.
    https://gocv.io/getting-started/linux/
    
#### Test locally

1. Open the assignment directory and execute `go run server.go`
2. open another terminal and cd to same directory.
3. Run `go run client.go`


