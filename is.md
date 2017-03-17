# Chotot Image Service ___CIS___
## Goal:
  Single, centralized image service, to power Chotot application

## Technology:
  - ___Ceph___ battle tested, distributed Object storage
  - ___Imaginary___ Golang based, image processing service, based on LibUV for fast image transformation. Interfacing on top of _Ceph_ Features included:
    - Resize
    - Crop
    - Watermark
    - Affine transformation
  - ___Varnish___ caching: Acting as delivery layer service to application.

## Roadmap:

## API

### Form Data:

If you're pushing images to `imaginary` as `multipart/form-data` (you can do it as well as `--data-binary`), you must define at least one input field ca
lled `file` with the raw image data in order to be processed properly by imaginary.

### Upload Image

#### POST /upload
Accepts: `multipart/form-data`
Support: `--data-binary` - Directly POST raw image data

#### Allowed variables
- service `string` `required` - Pool (Ceph concept) or Bucket (S3 concept)
- oid `string` `required` - Image ID

####  `GET <hash_key>/<service>/<oid>/<action>/`

EX:
```
GET /6e3df04527daf6a25989aa9345b47bc7d644ea18/property_project/1245 - Full size for profile picture
GET /6e3df04527daf6a25989aa9345b47bc7d644ea18/property_project/1245/thumbnail?width=100 - Thumbnail image
```

### Fetching Image Directly

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
