# Rainbow 🌈

Rainbow is a tools written in golang working for Mac.

It auto set current earth's photo from [http://himawari8.nict.go.jp/](http://himawari8.nict.go.jp/) as your desktop wallpaper. And can also download the past photos of earth.

## For Mac OSX

Mac ui wrap see [https://github.com/EyciaZhou/rainbow-osx-swift](https://github.com/EyciaZhou/rainbow-osx-swift)

## Requirements

- github.com/everdev/mack
- github.com/EyciaZhou/geo.go

## Build

    go build -o rainbow-cli


## Usage

	Usage of ./rainbow-cli:
	  -ang float
	    	Rotation angle, expressed in degrees, eg. 60, and a float is accpetable
	  -d int
	    	the degree of the picture, with higher value of 'd', the final picutre with higher quality, but use more network flow. 'd' must be the power of 2, like 1, 2, 4, 8, and at most 16 (default 4)
	  -f	force download, remove the recent picture if in need
	  -yst
	    	download all pictures of yesterday, if with this argument, this command will ignore other arguments
