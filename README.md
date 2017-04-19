
## Introduction
This is a fork repo from [hn2on](https://github.com/h2non/imaginary) that is built to support integration of Ceph Object Storage.  
Please check [hn2on](https://github.com/h2non/imaginary)'s repo to see original README.

## Supported image operations

- Resize
- Enlarge
- Crop
- Rotate (with auto-rotate based on EXIF orientation)
- Flip (with auto-flip based on EXIF metadata)
- Flop
- Zoom
- Thumbnail
- Configurable image area extraction
- Embed/Extend image, supporting multiple modes (white, black, mirror, copy or custom background color)
- Watermark (customizable by text)
- Custom output color space (RGB, black/white...)
- Format conversion (with additional quality/compression settings)
- Info (image size, format, orientation, alpha...)
- Reply with default or custom placeholder image in case of error.

## Prerequisites

- [libvips](https://github.com/jcupitt/libvips) v7.40.0+ or 8+ (8.3+ recommended)
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- Go 1.3+

## Installation

```bash
go get -u github.com/hieutrtr/imaginary
```

Also, be sure you have the latest version of `bimg`:
```bash
go get -u gopkg.in/h2non/bimg.v1
```

### libvips

Run the following script as `sudo` (supports OSX, Debian/Ubuntu, Redhat, Fedora, Amazon Linux):
```bash
curl -s https://raw.githubusercontent.com/h2non/bimg/master/preinstall.sh | sudo bash -
```

The [install script](https://github.com/h2non/bimg/blob/master/preinstall.sh) requires `curl` and `pkg-config`

### librados
Install on Debian/Ubuntu distributions.
```
sudo apt-get install librados-dev
```
Install on RHEL/CentOS distributions.
```
sudo yum install librados2-devel
```

### Docker

See [Dockerfile](https://github.com/hieutrtr/imaginary/blob/master/Dockerfile) for image details.

Fetch the image (comes with latest stable Go and libvips versions)
```
docker pull registry.gitlab.com/hieutrtr/imaginary
```

Start the container with optional flags (default listening on port 9000)
```
docker run -p 9000:9000 registry.gitlab.com/hieutrtr/imaginary -cors -gzip
```

Start the container in debug mode:
```
docker run -p 9000:9000 -e "DEBUG=*" registry.gitlab.com/hieutrtr/imaginary
```

Enter to the interactive shell in a running container
```
sudo docker exec -it <containerIdOrName> bash
```

Stop the container
```
docker stop hieutrtr/imaginary
```

You can see all the Docker tags [here](https://hub.docker.com/r/hieutrtr/imaginary/tags/).

## Performance

libvips is probably the faster open source solution for image processing.
Here you can see some performance test comparisons for multiple scenarios:

- [libvips speed and memory usage](http://www.vips.ecs.soton.ac.uk/index.php?title=Speed_and_Memory_Use)
- [bimg](https://github.com/h2non/bimg#Performance) (Go library with C bindings to libvips)

## Benchmark

See [benchmark.sh](https://github.com/h2non/imaginary/blob/master/benchmark.sh) for more details

Environment: Go 1.4.2. libvips-7.42.3. OSX i7 2.7Ghz

```
Requests  [total]       200
Duration  [total, attack, wait]   10.030639787s, 9.949499515s, 81.140272ms
Latencies [mean, 50, 95, 99, max]   83.124471ms, 82.899435ms, 88.948008ms, 95.547765ms, 104.384977ms
Bytes In  [total, mean]     23443800, 117219.00
Bytes Out [total, mean]     175517000, 877585.00
Success   [ratio]       100.00%
Status Codes  [code:count]      200:200
```

## Command-line usage

```
Usage:
  imaginary -p 80
  imaginary -cors -gzip
  imaginary -concurrency 10
  imaginary -path-prefix /api/v1
  imaginary -enable-url-source
  imaginary -enable-url-source -allowed-origins http://localhost,http://server.com
  imaginary -enable-url-source -enable-auth-forwarding
  imaginary -enable-url-source -authorization "Basic AwDJdL2DbwrD=="
	imaginary -enable-placeholder
	imaginery -enable-url-source -placeholder ./placeholder.jpg
	imaginary -enable-ceph -ceph-config /etc/ceph/ceph.conf
	imaginary -enable-safe-route -safe-key "secret-hash"
	imaginary -enable-tracking
	imaginary -h | -help
  imaginary -v | -version

Options:
  -a <addr>                 bind address [default: *]
  -p <port>                 bind port [default: 8088]
  -h, -help                 output help
  -v, -version              output version
  -path-prefix <value>      Url path prefix to listen to [default: "/"]
  -cors                     Enable CORS support [default: false]
  -gzip                     Enable gzip compression [default: false]
  -key <key>                Define API key for authorization
  -mount <path>             Mount server local directory
  -http-cache-ttl <num>     The TTL in seconds. Adds caching headers to locally served files.
  -http-read-timeout <num>  HTTP read timeout in seconds [default: 30]
  -http-write-timeout <num> HTTP write timeout in seconds [default: 30]
  -enable-url-source        Restrict remote image source processing to certain origins (separated by commas)
	-enable-placeholder       Enable image response placeholder to be used in case of error [default: false]
  -enable-auth-forwarding   Forwards X-Forward-Authorization or Authorization header to the image source server. -enable-url-source flag must be defined. Tip: secure your server from public access to prevent attack vectors
  -allowed-origins <urls>   TLS certificate file path
  -certfile <path>          TLS certificate file path
  -keyfile <path>           TLS private key file path
  -authorization <value>    Defines a constant Authorization header value passed to all the image source servers. -enable-url-source flag must be defined. This overwrites authorization headers forwarding behavior via X-Forward-Authorization
  -placeholder <path>       Image path to image custom placeholder to be used in case of error. Recommended minimum image size is: 1200x1200
	-concurrency <num>         Throttle concurrency limit per second [default: disabled]
  -burst <num>              Throttle burst max cache size [default: 100]
  -mrelease <num>           OS memory release interval in seconds [default: 30]
  -cpus <num>               Number of used cpu cores.
                            (default for current machine is 8 cores)

	-enable-ceph              enable ceph integration
	-ceph-config              path to ceph config
	-enable-safe-route				enable safe route url
	-safe-key									secret key to hash URI that is used with enable-safe-route
	-enable-tracking 					tracking event
```

Start script for Ceph integration:
```bash
export KAFKA_BROKERS=kafka1:9092,kafka2:9092,kafka3:9092
imaginary -http-cache-ttl 86400 -enable-url-source -concurrency 100 -enable-ceph -enable-safe-route -safe-key randomkey -enable-tracking -p 8088
```
enable-safe-route require safe-key: in this case is `randomkey`
For generating token : hmac(randomkey,sha1(path))
```
GET /project/image_id/thumbnail?width=100
token = hmac(random, sha1('/project/image_id/thumbnail?width=100'))
GET /project/image_id
token = hmac(random, sha1('/project/image_id'))
```

#### Examples

Uploading a local image:
```
curl -XPOST "http://localhost:8088/upload/project/image_id"
              --data-binary @"/home/user1/Desktop/test.jpg"
```

Fetching the image from Ceph:
```
curl "http://localhost:8088/7d6e3468f88cccbde9e94062650d632786dd54ea/project/image_id/thumbnail?width=100"
```

### Params

Complete list of available params. Take a look to each specific endpoint to see which params are supported.
Image measures are always in pixels, unless otherwise indicated.

- **width**       `int`   - Width of image area to extract/resize
- **height**      `int`   - Height of image area to extract/resize
- **top**         `int`   - Top edge of area to extract. Example: `100`
- **left**        `int`   - Left edge of area to extract. Example: `100`
- **areawidth**   `int`   - Height area to extract. Example: `300`
- **areaheight**  `int`   - Width area to extract. Example: `300`
- **quality**     `int`   - JPEG image quality between 1-100. Defaults to `80`
- **compression** `int`   - PNG compression level. Default: `6`
- **rotate**      `int`   - Image rotation angle. Must be multiple of `90`. Example: `180`
- **factor**      `int`   - Zoom factor level. Example: `2`
- **margin**      `int`   - Text area margin for watermark. Example: `50`
- **dpi**         `int`   - DPI value for watermark. Example: `150`
- **textwidth**   `int`   - Text area width for watermark. Example: `200`
- **opacity**     `float` - Opacity level for watermark text. Default: `0.2`
- **flip**        `bool`  - Transform the resultant image with flip operation. Default: `false`
- **flop**        `bool`  - Transform the resultant image with flop operation. Default: `false`
- **force**       `bool`  - Force image transformation size. Default: `false`
- **nocrop**      `bool`  - Disable crop transformation enabled by default by some operations. Default: `false`
- **noreplicate** `bool`  - Disable text replication in watermark. Defaults to `false`
- **norotation**  `bool`  - Disable auto rotation based on EXIF orientation. Defaults to `false`
- **noprofile**   `bool`  - Disable adding ICC profile metadata. Defaults to `false`
- **text**        `string` - Watermark text content. Example: `copyright (c) 2189`
- **font**        `string` - Watermark text font type and format. Example: `sans bold 12`
- **color**       `string` - Watermark text RGB decimal base color. Example: `255,200,150`
- **type**        `string` - Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`
- **gravity**     `string` - Define the crop operation gravity. Supported values are: `north`, `south`, `centre`, `west` and `east`. Defaults to `centre`.
- **file**        `string` - Use image from server local file path. In order to use this you must pass the `-mount=<dir>` flag.
- **url**         `string` - Fetch the image from a remove HTTP server. In order to use this you must pass the `-enable-url-source` flag.
- **colorspace**  `string` - Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)
- **field**       `string` - Custom image form field name if using `multipart/form`. Defaults to: `file`
- **extend**      `string` - Extend represents the image extend mode used when the edges of an image are extended. Allowed values are: `black`, `copy`, `mirror`, `white` and `background`. If `background` value is specified, you can define the desired extend RGB color via `background` param, such as `?extend=background&background=250,20,10`. For more info, see [libvips docs](http://www.vips.ecs.soton.ac.uk/supported/8.4/doc/html/libvips/libvips-conversion.html#VIPS-EXTEND-BACKGROUND:CAPS).
- **background**  `string` - Background RGB decimal base color to use when flattening transparent PNGs. Example: `255,200,150`

#### GET /
Content-Type: `application/json`

Serves as JSON the current `imaginary`, `bimg` and `libvips` versions.

Example response:
```json
{
  "imaginary": "0.1.28",
  "bimg": "1.0.5",
  "libvips": "8.4.1"
}
```

#### GET /health
Content-Type: `application/json`

Provides some useful statistics about the server stats with the following structure:

- **uptime** `number` - Server process uptime in seconds.
- **allocatedMemory** `number` - Currently allocated memory in megabytes.
- **totalAllocatedMemory** `number` - Total allocated memory over the time in megabytes.
- **gorouting** `number` - Number of running gorouting.
- **cpus** `number` - Number of used CPU cores.

Example response:
```json
{
  "uptime": 1293,
  "allocatedMemory": 5.31,
  "totalAllocatedMemory": 34.3,
  "goroutines": 19,
  "cpus": 8
}
```

#### GET /form
Content Type: `text/html`

Serves an ugly HTML form, just for testing/playground purposes

#### GET | POST /info
Accepts: `image/*, multipart/form-data`. Content-Type: `application/json`

Returns the image metadata as JSON:
```json
{
  "width": 550,
  "height": 740,
  "type": "jpeg",
  "space": "srgb",
  "hasAlpha": false,
  "hasProfile": true,
  "channels": 3,
  "orientation": 1
}
```

#### GET | POST /crop
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

Crop the image by a given width or height. Image ratio is maintained

##### Allowed params

- width `int`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- embed `bool`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- gravity `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /resize
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

Resize an image by width or height. Image aspect ratio is maintained

##### Allowed params

- width `int` `required`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /enlarge
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- width `int` `required`
- height `int` `required`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /extract
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- top `int` `required`
- left `int`
- areawidth `int` `required`
- areaheight `int`
- width `int`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /zoom
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- factor `number` `required`
- width `int`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /thumbnail
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- width `int`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /rotate
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- rotate `int` `required`
- width `int`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /flip
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- width `int`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /flop
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- width `int`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /convert
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- type `string` `required`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /watermark
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*`

##### Allowed params

- text `string` `required`
- margin `int`
- dpi `int`
- textwidth `int`
- opacity `float`
- noreplicate `bool`
- font `string`
- color `string`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- embed `bool`
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- flip `bool`
- flop `bool`
- extend `string`
- background `string` - Example: `?background=250,20,10`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads
