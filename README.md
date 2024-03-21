
# Chalanges

1. Convert image files from PNG to JPEG.
2. Resize images according to specified dimensions.
3. Compress images to reduce file size while maintaining reasonable quality.

## My Assumtion
1. Since the endpoints would process multiple file, response would be send in `zip` format.
2. There is a three endpoints that answer above funtionalities, and one endpoint that combine all functionalities. Postman documenttaion is provided.

# How to Run The Server

## Run Manual
To run manual localy, install dependency first. This project use OpenCV and [gocv](https://github.com/hybridgroup/gocv) for image processing. I use this [installation](https://github.com/hybridgroup/gocv?tab=readme-ov-file#how-to-install) instruction to install on my local machine. 

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


## Run using Docker
No need to install dependency if you run using docker

Build docker localy:
```
make docker-build
```

run docker container:
```
docker run -d image-processing-go
```


# How to run unit tests

```
make test
```

to view html coverage
```
make test-view-html
```

# API Docs
Postman API Documentation is provided in [docs](docs)

## End-point: Png to Jpeg
### Method: POST
>```
>{{SERVER}}/png-to-jpeg
>```
### Body formdata

|Param|value|Type|
|---|---|---|
|files[]|/dir/subdir/flower.png|file|
|files[]|/dir/.cache/car-967387_1920.png|file|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Compress
### Method: POST
>```
>{{SERVER}}/compress
>```
### Body formdata

|Param|value|Type|
|---|---|---|
|files[]|/dir/subdir/car-967387_1920.png|file|
|files[]|/dir/subdir/cat.jpg|file|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Resize
### Method: POST
>```
>{{SERVER}}/resize
>```
### Body formdata

|Param|value|Type|
|---|---|---|
|height[]|90|text|
|height[]|90|text|
|width[]|90|text|
|width[]|90|text|
|files[]|/dir/subdir/cat.jpg|file|
|files[]|/dir/subdir/car-967387_1920.png|file|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Process Image
Combine 3 Functionalities (Convert png to jpeg, resize, and compress)
### Method: POST
>```
>{{SERVER}}http://localhost:10000
>```
### Body formdata

|Param|value|Type|
|---|---|---|
|width[]|50|text|
|width[]|50|text|
|height[]|80|text|
|height[]|90|text|
|files[]|/dir/subdir/flower.png|file|
|files[]|/dir/.cache/car-967387_1920.png|file|
