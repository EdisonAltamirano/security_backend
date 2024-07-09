<img src="readme_resources/nawilogo.jpg" width="100" ><img> 

# Nawi backend
This repository contains the development of camera streamer backend which is the software project in charge of sending the live transmission of a TCP/IP camera, it has to be connected to nawi camera service with the correct IP of the container.

## Project setup

1. Clone the project repository on your local machine.

   SSH:

   ```bash
   $ git clone --recurse-submodules https://github.com/RoBorregos/Robocup-Home.git
   ```
2.  Build the image just the first time
  ```bash
   $ docker build -t backend:nawi -f Dockerfile .
  ```
3. Run the container containing the code 
  ```bash
   $ docker run --rm -it -p 3002:3000/tcp backend:nawi
  ```
6. Enter the container and inside security_backend/cmd/security_backend, run the following code.
  ```bash
   $ go run .
  ```
7. Consider that the port has to be 3002 and the camera_service has to be the IP con the container running camera_service
  ```bash
   $ port = 3002
    [camera_service]
    hostname = '172.17.0.5'
    port = 3001
  ```
## Useful commands
1. The camera link is generated in https://www.ispyconnect.com/camera/imou
```bash
rtsp://admin:L2FDAF98@192.168.1.73:554/cam/realmonitor?channel=1&subtype=0&unicast=true&proto=Onvif
```
