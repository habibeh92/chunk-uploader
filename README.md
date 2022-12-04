
# ASCII Art Chunk Uploader
A simple webservice to upload chunked ASCII art images and download them as a single image.

## How to Run
You can easily serve the webservice on port `4444` via docker-compose using command below:
```
docker-compose up --build
```
Then you can access the `http://localhost:4444` to register your image, upload the chunks and download the file as `text/plain`

## Available Endpoints

+ **Registering an image**:
    + **method**: `POST`
    + **URI**: `/image`
    + **Content-Type**: `application/json`
    + **Request Body**:

        ```json
        {
          "sha256": "abc123easyasdoremi...",
          "size": 123456,
          "chunk_size": 256
        }
        ```

    + **Responses**:
  
      | Code                       |              Description           |
      |----------------------------|------------------------------------|
      | 201 Created                | Image successfully registered       |
      | 409 Conflict               | Image already exists               |
      | 400 Bad Request            | Malformed request                  |
      | 415 Unsupported Media Type | Unsupported payload format         |

+ **Uploading an image chunk**:
    + **method**: `POST`
    + **URI**: `/image/<sha256>/chunks`
    + **Content-Type**: `application/json`
    + **Request Body**:

        ```json
        {
          "id": 1,
          "size": 256,
          "data": "8   888   , 888    Y888 888 888    ,ee 888 888 888 888 ...",
        }
        ```

    + **Responses**:
  
      | Code          |              Description           |
      |---------------|------------------------------------|
      | 201 Created   | Chunk successfully uploaded         |
      | 409 Conflict  | Chunk already exists               |
      | 404 Not Found | Image not found                    |

+ **Downloading an image**:
    + **method**: `GET`
    + **URI**: `/image/<sha256>`
    + **Accept**: `text/plain`
    + **Responses**:
  
      | Code          |              Description           |
      |---------------|------------------------------------|
      | 200 OK        | Image successfully downloaded       |
      | 404 Not Found | Image not found                    |

+ **Errors**:
    + **Accept**: `application/json`
    + **Response body**:

      ```json
      {
        "code": "400",
        "message": "Chunk ID field is missing."
      }
      ```
