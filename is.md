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

If you're pushing images to `imaginary` as `multipart/form-data` (you can do it as well as `--data-binary`), you must define at least one input field called `file` with the raw image data in order to be processed properly by imaginary.

### Upload Image

#### POST /upload
Accepts: `multipart/form-data`
Support: `--data-binary` - Directly POST raw image data

#### Allowed params
- cpool `string` `required` - Pool (Ceph concept) or Namespace of service
- coid `string` `required` - Image ID

### Fetching Image Through Varnish
#####  `GET /<service>_<action>/<img_id>`

EX:
`GET /profile/190790 - Full size for profile picture`
`GET /property_project_thumb/190790 - Thumbnail image`
`GET /ads_wm/190790 - Watermark image`

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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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

- cpool `string` - Specific Ceph Pool (for every service)
- coid `string` - Specific object id of the Pool of Service
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
