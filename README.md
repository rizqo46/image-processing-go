
# Chalanges

1. Convert image files from PNG to JPEG.
2. Resize images according to specified dimensions.
3. Compress images to reduce file size while maintaining reasonable quality.

## My Assumtion
1. Since the endpoints would process multiple file, response would be send in `zip` format.
2. There is a three endpoints that answer above funtionalities, and one endpoint that combine all functionalities. Postman documenttaion is provided.

# How to Run The Server

## Run Manual
This project use OpenCV and [gocv](https://github.com/hybridgroup/gocv) for image processing. I use this [installation](https://github.com/hybridgroup/gocv?tab=readme-ov-file#how-to-install) instruction to install on my local machine. 

**Note**: If you run on linux machine and getting error on installation you may need to edit the makefile, check this [issue](https://github.com/hybridgroup/gocv/issues/978) on github fo the instruction.

to run project
```
make run
```

run using build
```
make build-and-run 
```
the server will run on port 8080 by default, export env PORT to run in specific port.


# How to run unit tests

```
make test
```

to view html coverage
```
make test-view-html
```