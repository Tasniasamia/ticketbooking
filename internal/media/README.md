# Media Module (Cloudinary + Central Media Table)

## Architecture

একটাই `media` টেবিল সব ধরনের ফাইল (image / video / pdf / file) রাখে।

| Column        | Description                                      |
|---------------|--------------------------------------------------|
| `id`          | Auto-generated DB primary key (এটা দিয়ে delete)  |
| `image_id`    | Cloudinary `public_id` (Cloudinary থেকে delete)  |
| `public_url`  | HTTP URL                                         |
| `secure_url`  | HTTPS URL                                        |
| `type`        | `image` \| `video` \| `pdf` \| `file`            |
| `model_name`  | `event` \| `user` \| `product` \| ...            |
| `model_id`    | Optional related entity id                       |
| `created_at`  |                                                  |
| `deleted_at`  | Soft delete                                      |

## Frontend থেকে কী পাঠাতে হবে (Upload)

`POST /api/v1/media/upload`  
`Content-Type: multipart/form-data`

| Field        | Required | Example          |
|--------------|----------|------------------|
| `file`       | ✅       | binary file      |
| `model_name` | ✅       | `event`          |
| `model_id`   | ❌       | `12`             |
| `folder`     | ❌       | `events/banners` |

### Example (curl)

```bash
curl -X POST http://localhost:1323/api/v1/media/upload \
  -F "file=@/path/to/poster.jpg" \
  -F "model_name=event" \
  -F "model_id=5"
```

### Success Response

```json
{
  "success": true,
  "statusCode": 201,
  "message": "File uploaded successfully",
  "data": {
    "media": {
      "id": 15,
      "image_id": "event/1723..._poster",
      "public_url": "http://res.cloudinary.com/...",
      "secure_url": "https://res.cloudinary.com/...",
      "type": "image",
      "model_name": "event",
      "model_id": 5,
      "created_at": "2026-08-08 11:30:00"
    },
    "keyed": {
      "event_url": "https://res.cloudinary.com/...",
      "event_id": 15,
      "image_id": "event/1723..._poster",
      "public_url": "http://res.cloudinary.com/...",
      "type": "image",
      "model_name": "event",
      "created_at": "2026-08-08 11:30:00"
    }
  }
}
```

`event_id` = media টেবিলের auto id। পরে delete করতে এই id ব্যবহার করবে।

## Delete

`DELETE /api/v1/media/:id`  
(`:id` = media table-এর auto generated id)

Flow:
1. DB থেকে row খুঁজে বের করে
2. `image_id` (Cloudinary public_id) দিয়ে Cloudinary থেকে destroy
3. DB-তে soft delete (`deleted_at` set)

## List / Filter

```
GET /api/v1/media?model_name=event&model_id=5&type=image&page=1&limit=10
```

## Get single

```
GET /api/v1/media/15
```

## Environment

`.env`-এ যোগ করো:

```
CLOUDINARY_CLOUD_NAME=xxx
CLOUDINARY_API_KEY=xxx
CLOUDINARY_API_SECRET=xxx
```

## Product / Event model-এ কীভাবে ব্যবহার করবে

Event create করার সময় image upload আগে করো → `event_id` (media id) পাবে → তারপর Event create/update-এ সেই media id রাখো।

অথবা Event table-এ `media_ids` JSONB array রাখতে পারো, অথবা আলাদা `event_media` pivot table।

GetEvent() করার সময় media id দিয়ে media table থেকে secure_url নিয়ে response-এ `images: [...]` পাঠাবে।
