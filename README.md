# img
pure go image rendering server. accepts a POST, stores the file on s3, and performs various manipulations on the images via GET string params

**example upload**

```
curl -i -F file=@photo2.jpg http://localhost:8080/upload

curl -i -F file=@photo2.jpg http://img.tinyfactory.io/upload
```

## api

- `action`: the action to perform. currently supports:
	crop: currently only cropping form the middle is supported
- `w`: width in pixels of the cropped image
- `h`: height in pixels of the cropped image

**example request**

```
http://img.tinyfactory.io/img/671e3346e405b99441bf4f0de7abc4dd?action=thumbnail&w=500&h=500
```